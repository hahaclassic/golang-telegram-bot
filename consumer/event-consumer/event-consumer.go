package event_consumer

import (
	"log"
	"time"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, bath int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: bath,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {

		if err := c.processor.Process(event); err != nil {
			log.Print(errhandling.Wrap("can't handle event", err))

			continue
		}
	}

	return nil
}
