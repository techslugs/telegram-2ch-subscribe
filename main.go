package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"telegram-2ch-news-bot/env_config"
)

func main() {
	config := env_config.ReadConfig()

	api, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := api.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}

	handleUpdates(api, updates)
}

func handleUpdates(api *tgbotapi.BotAPI, updates <-chan tgbotapi.Update) {
	for update := range updates {
		if update.Message == nil {
			log.Printf("%v", update)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		api.Send(msg)
	}
}
