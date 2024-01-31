package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	tgClient "github.com/hahaclassic/golang-telegram-bot.git/clients/telegram"
	event_consumer "github.com/hahaclassic/golang-telegram-bot.git/consumer/event-consumer"
	"github.com/hahaclassic/golang-telegram-bot.git/events/telegram"
	"github.com/hahaclassic/golang-telegram-bot.git/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/data.db"
	batchSize         = 100
	errChanSize       = 100
)

var (
	mainToken   string
	loggerToken string
	adminChatID int
)

func init() {
	flag.StringVar(&mainToken, "tg-bot-token", "", "token for access to main telegram bot")
	flag.StringVar(&loggerToken, "logger", "", "token for access to main logger bot")
	flag.IntVar(&adminChatID, "admin", -1, "admin chatID")
	flag.Parse()

	if mainToken == "" || loggerToken == "" || adminChatID == -1 {
		fmt.Println(mainToken, loggerToken, adminChatID)
		log.Fatal("token is not specified")
	}
}

func main() {
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		err := os.MkdirAll("./data/sqlite/", 0777)

		if err != nil {
			log.Fatal("can't create directory")
		}
	}

	// Create database
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: %s", err)
	}

	err = s.Init(context.TODO())
	if err != nil {
		log.Fatalf("can't init storage: %s", err)
	}

	// Create events Processor
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mainToken),
		tgClient.New(tgBotHost, loggerToken),
		adminChatID,
		s,
	)

	log.Print("[START]")

	// Create consumer
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize, errChanSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
