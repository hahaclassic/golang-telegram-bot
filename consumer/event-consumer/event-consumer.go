package event_consumer

import (
	"log"
	"sync"
	"time"

	"github.com/hahaclassic/golang-telegram-bot.git/events"
	"github.com/hahaclassic/golang-telegram-bot.git/lib/errhandling"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
	chanSize  int
}

func New(fetcher events.Fetcher, processor events.Processor, batсh int, chanSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batсh,
		chanSize:  chanSize,
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

		go c.handleEvents(gotEvents)
	}
}

func (c *Consumer) handleEvents(events []events.Event) {
	errors := make(chan error, c.chanSize)

	wg := sync.WaitGroup{}
	wg.Add(len(events))

	for _, event := range events {

		go c.processor.Process(event, errors, &wg)

		if err := <-errors; err != nil {
			log.Print(errhandling.Wrap("can't handle event", err))

			continue
		}
		// if err := c.processor.Process(event); err != nil {
		// 	log.Print(errhandling.Wrap("can't handle event", err))

		// 	continue
		// }
	}
	wg.Wait()
}
