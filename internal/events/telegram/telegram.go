package telegram

import (
	"context"
	"errors"
	"sync"

	tgclient "github.com/hahaclassic/golang-telegram-bot.git/internal/clients/telegram"
	"github.com/hahaclassic/golang-telegram-bot.git/internal/events"
	"github.com/hahaclassic/golang-telegram-bot.git/internal/storage"
	"github.com/hahaclassic/golang-telegram-bot.git/internal/storage/sqlite"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

// Данный тип реализует сразу два интерфейса: Processor() и Fetcher()
type Processor struct {
	tg          *tgclient.Client
	logger      *tgclient.Client
	adminChatID int
	offset      int
	storage     storage.Storage
	sessions    map[int]*Session
}

// Убрать status (совместить с currOperation)
type Session struct {
	currentOperation Operation
	url              string
	tag              string
	folderID         string
	lastMessageID    int
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
	ErrNoFolders       = errors.New("no existing folders")
	ErrEmptyFolder     = errors.New("empty folder")
)

func New(config *Config) (*Processor, error) {

	s, err := sqlite.New(config.StoragePath)
	if err != nil {
		return nil, err
	}

	err = s.Init(context.TODO())
	if err != nil {
		return nil, err
	}

	processor := &Processor{
		tg:          tgclient.New(config.Host, config.MainToken),
		logger:      tgclient.New(config.Host, config.LoggerToken),
		adminChatID: config.AdminChatID,
		storage:     s,
		sessions:    make(map[int]*Session),
	}

	return processor, nil
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
		if err != nil || p.sessions[event.UserID].currentOperation == DoneCmd {
			delete(p.sessions, event.UserID)
		}
	}()

	if _, ok := p.sessions[event.UserID]; !ok {
		p.sessions[event.UserID] = &Session{currentOperation: DoneCmd}
	}

	return p.doCallbackCmd(&event, meta)
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't process message", err) }()

	defer func() {
		if err != nil || p.sessions[event.UserID].currentOperation == DoneCmd {
			delete(p.sessions, event.UserID)
		}
	}()

	if _, ok := p.sessions[event.UserID]; !ok {
		p.sessions[event.UserID] = &Session{currentOperation: DoneCmd}
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

func event(upd tgclient.Update) events.Event {

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
		res.Meta = CallbackMeta{
			Data:    upd.CallbackQuery.Data,
			QueryID: upd.CallbackQuery.QueryID,
		}
	}

	return res
}

func fetchType(upd tgclient.Update) events.Type {
	switch {
	case upd.Message != nil:
		return events.Message
	case upd.CallbackQuery != nil:
		return events.CallbackQuery
	default:
		return events.Unknown
	}
}
