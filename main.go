package main

import (
	"gopkg.in/mgo.v2"
	"log"
	"telegram-2ch-news-bot/bot"
	"telegram-2ch-news-bot/env_config"
	"telegram-2ch-news-bot/storage"
	"time"
)

func setupStorage(config env_config.Config) (*storage.Storage, error) {
	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		return nil, err
	}

	storage, err := storage.NewStorage(session.DB(config.MongoDatabase))
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func main() {
	config := env_config.ReadConfig()

	storage, err := setupStorage(config)
	if err != nil {
		log.Panic(err)
	}

	err = bot.StartBot(
		config.TelegramToken,
		time.Second*time.Duration(config.BoardsListPollingTimeout),
		time.Second*time.Duration(config.BoardPollingTimeout),
		storage,
	)

	if err != nil {
		log.Panic(err)
	}
	select {}
}
