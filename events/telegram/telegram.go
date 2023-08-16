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
	tg      *tgClient.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrUnknownEvent    = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *tgClient.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
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
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errhandling.Wrap("can't process the message", ErrUnknownEvent)
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = errhandling.WrapIfErr("can't process message", err) }()

	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
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

func event(upd tgClient.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd tgClient.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd tgClient.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
