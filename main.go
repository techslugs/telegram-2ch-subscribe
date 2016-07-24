package main

import (
	"log"
	"github.com/tmwh/telegram-2ch-subscribe/bot"
	"github.com/tmwh/telegram-2ch-subscribe/env_config"
	"github.com/tmwh/telegram-2ch-subscribe/storage"
	"github.com/tmwh/telegram-2ch-subscribe/telegram"
	"github.com/tmwh/telegram-2ch-subscribe/web"
	"time"
)

func main() {
	config := env_config.ReadConfig()

	storage, err := storage.NewStorage(config.MongoURL, config.MongoDatabase)
	if err != nil {
		log.Panic(err)
	}

	telegramClient, err := telegram.NewClient(config.TelegramToken, storage)
	if err != nil {
		log.Panic(err)
	}

	err = bot.StartBot(
		time.Second*time.Duration(config.BoardsListPollingTimeout),
		time.Second*time.Duration(config.BoardPollingTimeout),
		telegramClient,
		storage,
	)
	if err != nil {
		log.Panic(err)
	}

	log.Fatal(web.StartServer(config.IpAddress, config.Port))

	select {}
}
