package main

import (
	"log"
	"telegram-2ch-news-bot/bot"
	"telegram-2ch-news-bot/env_config"
	"time"
)

func main() {
	config := env_config.ReadConfig()
	err := bot.StartBot(config.TelegramToken, time.Second*time.Duration(config.BoardPollingTimeout))
	if err != nil {
		log.Panic(err)
	}
	select {}
}
