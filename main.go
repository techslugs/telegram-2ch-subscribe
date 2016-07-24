package main

import (
	"log"
	"telegram-2ch-news-bot/bot"
	"telegram-2ch-news-bot/env_config"
	"telegram-2ch-news-bot/storage"
	"telegram-2ch-news-bot/telegram"
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
	select {}
}
