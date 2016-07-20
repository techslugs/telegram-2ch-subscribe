package main

import (
	"log"
	"telegram-2ch-news-bot/bot"
	"telegram-2ch-news-bot/env_config"
)

func main() {
	config := env_config.ReadConfig()
	err := bot.StartBot(config.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	for {
	}
}
