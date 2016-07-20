package main

import (
	"gopkg.in/mgo.v2"
	"log"
	"telegram-2ch-news-bot/bot"
	"telegram-2ch-news-bot/env_config"
	"time"
)

func setupStorage(config env_config.Config) (*bot.Storage, error) {
	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		return nil, err
	}

	storage, err := bot.NewStorage(session.DB(config.MongoDatabase))
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

	err = storage.UpdateBoardTimestamp("b", time.Now().Add(-5*time.Minute).Unix())
	if err != nil {
		log.Panic(err)
	}

	err = bot.StartBot(
		config.TelegramToken,
		time.Second*time.Duration(config.BoardPollingTimeout),
		storage,
	)

	if err != nil {
		log.Panic(err)
	}
	select {}
}
