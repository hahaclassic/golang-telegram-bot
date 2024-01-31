package telegram

import (
	"errors"
	"sync"

	tgClient "github.com/hahaclassic/golang-telegram-bot.git/clients/telegram"
	"github.com/hahaclassic/golang-telegram-bot.git/events"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

// Данный тип реализует сразу два интерфейса: Processor() и Fetcher()
type Processor struct {
	tg          *tgClient.Client
	logger      *tgClient.Client
	adminChatID int
	offset      int
	storage     storage.Storage
	sessions    map[int]*Session
}

// При статусе ОК сессия будет удалятся из карты. (Возможно, будет удаляется через некоторое время)
type Session struct {
	currentOperation string
	url              string
	name             string
	folder           string
	status           bool
}

type MessageMeta struct {
	ChatID int
	UserID int
}

type CallbackMeta struct {
	QueryID string
	UserID  int
	Message string
	ChatID  int
}

const (
	statusOK         = true
	statusProcessing = false
)

var (
	ErrUnknownEvent    = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
	ErrNoFolders       = errors.New("No existing folders")
	ErrEmptyFolder     = errors.New("Empty folder")
)

func New(client *tgClient.Client, logger *tgClient.Client, adminChatID int, storage storage.Storage) *Processor {
	return &Processor{
		tg:          client,
		logger:      logger,
		adminChatID: adminChatID,
		storage:     storage,
		sessions:    make(map[int]*Session),
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errhandling.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event, errors chan error, wg *sync.WaitGroup) {

	defer wg.Done()

	switch event.Type {
	case events.Message:
		errors <- p.processMessage(event)
	case events.CallbackQuery:
		errors <- p.processCallbackQuery(event)
	default:
		errors <- errhandling.Wrap("can't process the message", ErrUnknownEvent)
	}
}

func (p *Processor) processCallbackQuery(event events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't process callback", err) }()

	meta, err := getCallbackMeta(event)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil || p.sessions[meta.UserID].status == statusOK {
			delete(p.sessions, meta.UserID)
		}
	}()

	if _, ok := p.sessions[meta.UserID]; !ok {
		p.sessions[meta.UserID] = &Session{
			status: statusOK,
		}
	}

	return p.doCallbackCmd(event.Text, &meta)
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't process message", err) }()

	meta, err := getMessageMeta(event)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil || p.sessions[meta.UserID].status == statusOK {
			delete(p.sessions, meta.UserID)
		}
	}()

	if _, ok := p.sessions[meta.UserID]; !ok {
		p.sessions[meta.UserID] = &Session{
			status: statusOK,
		}
		return p.startCmd(event.Text, meta.ChatID, meta.UserID)
	}

	return p.handleCmd(event.Text, meta.ChatID, meta.UserID)
}

func getMessageMeta(event events.Event) (MessageMeta, error) {
	res, ok := event.Meta.(MessageMeta)
	if !ok {
		return MessageMeta{}, errhandling.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func getCallbackMeta(event events.Event) (CallbackMeta, error) {
	res, ok := event.Meta.(CallbackMeta)
	if !ok {
		return CallbackMeta{}, errhandling.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd tgClient.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = MessageMeta{
			ChatID: upd.Message.Chat.ID,
			UserID: upd.Message.From.UserID,
		}
	} else if updType == events.CallbackQuery {
		res.Meta = CallbackMeta{
			QueryID: upd.CallbackQuery.QueryID,
			UserID:  upd.CallbackQuery.From.UserID,
			Message: upd.CallbackQuery.Message.Text,
			ChatID:  upd.CallbackQuery.Message.Chat.ID,
		}
	}

	return res
}

func fetchText(upd tgClient.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	} else if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}

	return ""
}

func fetchType(upd tgClient.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	} else if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}

	return events.Unknown
}
