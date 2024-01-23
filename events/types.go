package events

import "sync"

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event, errors chan error, wg *sync.WaitGroup)
}

type Type int

const (
	Unknown Type = iota
	Message
	CallbackQuery
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
