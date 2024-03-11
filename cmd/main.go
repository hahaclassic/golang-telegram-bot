package main

import (
	"log"

	consumer "github.com/hahaclassic/golang-telegram-bot.git/internal/consumer"
	"github.com/hahaclassic/golang-telegram-bot.git/internal/events/telegram"
)

const (
	configPath = "./configs/config.yaml"
)

func main() {

	config, err := telegram.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	eventsProcessor, err := telegram.New(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("[START]")

	// Create consumer
	cons := consumer.New(eventsProcessor, eventsProcessor, config.BatchSize, config.ErrChanSize)
	if err := cons.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
