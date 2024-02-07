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
	tag              string
	folderName       string
	status           bool
}

type CallbackMeta struct {
	Data    string
	QueryID string
}

const (
	statusOK         = true
	statusProcessing = false
	maxAttemts       = 100
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
		if err != nil || p.sessions[event.UserID].status == statusOK {
			delete(p.sessions, event.UserID)
		}
	}()

	if _, ok := p.sessions[event.UserID]; !ok {
		p.sessions[event.UserID] = &Session{status: statusOK}
	}

	return p.doCallbackCmd(&event, meta)
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't process message", err) }()

	defer func() {
		if err != nil || p.sessions[event.UserID].status == statusOK {
			delete(p.sessions, event.UserID)
		}
	}()

	if _, ok := p.sessions[event.UserID]; !ok {
		p.sessions[event.UserID] = &Session{status: statusOK}
		return p.startCmd(&event)
	}

	return p.handleCmd(&event)
}

func getCallbackMeta(event events.Event) (*CallbackMeta, error) {
	meta, ok := event.Meta.(CallbackMeta)
	if !ok {
		return nil, errhandling.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return &meta, nil
}

func event(upd tgClient.Update) events.Event {

	res := events.Event{
		Type: fetchType(upd),
	}

	if res.Type == events.Message {
		res.ChatID = upd.Message.Chat.ID
		res.Text = upd.Message.Text
		res.UserID = upd.Message.From.UserID
		res.Username = upd.Message.From.Username
	} else if res.Type == events.CallbackQuery {
		res.ChatID = upd.CallbackQuery.Message.Chat.ID
		res.Text = upd.CallbackQuery.Message.Text
		res.UserID = upd.CallbackQuery.From.UserID
		res.Username = upd.CallbackQuery.From.Username
	}

	return res
}

func fetchType(upd tgClient.Update) events.Type {
	switch {
	case upd.Message != nil:
		return events.Message
	case upd.CallbackQuery != nil:
		return events.CallbackQuery
	default:
		return events.Unknown
	}
}
