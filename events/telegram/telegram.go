package telegram

import (
	"errors"

	tgClient "github.com/hahaclassic/golang-telegram-bot.git/clients/telegram"
	"github.com/hahaclassic/golang-telegram-bot.git/events"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
	"github.com/hahaclassic/golang-telegram-bot.git/storage"
)

// Данный тип реализует сразу два интерфейса: Processor() и Fetcher()
type Processor struct {
	tg               *tgClient.Client
	offset           int
	currentOperation string
	lastMessage      string
	status           bool
	storage          storage.Storage
}

type Meta struct {
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
)

func New(client *tgClient.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		status:  statusOK,
		storage: storage,
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

func (p *Processor) Process(event events.Event) error {
	//log.Println(event)
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.CallbackQuery:
		return p.processCallbackQuery(event)
	default:
		return errhandling.Wrap("can't process the message", ErrUnknownEvent)
	}
}

func (p *Processor) processCallbackQuery(event events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't process callback", err) }()

	meta, err := callbackMeta(event)
	if err != nil {
		return err
	}

	if err := p.doCallbackCmd(event.Text, &meta); err != nil {
		return err
	}

	return nil
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't process message", err) }()

	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserID); err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errhandling.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func callbackMeta(event events.Event) (CallbackMeta, error) {
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
		res.Meta = Meta{
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
