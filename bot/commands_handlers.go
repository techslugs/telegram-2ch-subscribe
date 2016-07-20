package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (bot *Bot) handleCommandsFromTelegram() {
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
