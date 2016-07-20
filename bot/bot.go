package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Bot struct {
	telegramApi     *tgbotapi.BotAPI
	telegramUpdates <-chan tgbotapi.Update
	threadUpdates   <-chan ThreadUpdate
}

func StartBot(telegramToken string) error {
	bot := Bot{}
	err := bot.setupTelegramApi(telegramToken)
	if err != nil {
		return err
	}

	go bot.handleTelegramUpdates()
	return nil
}

func (bot *Bot) setupTelegramApi(telegramToken string) error {
	api, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	telegramUpdates, err := api.GetUpdatesChan(updateConfig)
	if err != nil {
		return err
	}

	bot.telegramApi = api
	bot.telegramUpdates = telegramUpdates
	return nil
}

func (bot *Bot) handleTelegramUpdates() {
	for update := range bot.telegramUpdates {
		if update.Message == nil {
			log.Printf("%v", update)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.telegramApi.Send(msg)
	}
}
