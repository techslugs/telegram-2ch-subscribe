package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tmwh/telegram-2ch-subscribe/bot/commands"
	"github.com/tmwh/telegram-2ch-subscribe/telegram"
	"log"
)

func StartHandleCommandsFromTelegram(telegramClient *telegram.Client) {
	for update := range telegramClient.TelegramUpdates {
		parseAndHandleCommand(telegramClient, &update)
	}
}

func parseAndHandleCommand(telegramClient *telegram.Client, update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	messageText := update.Message.Text
	switch {
	case commands.SubscribeChannel.Matches(messageText):
		commands.Handle(commands.SubscribeChannel, telegramClient, update.Message)
	case commands.UnsubscribeChannel.Matches(messageText):
		commands.Handle(commands.UnsubscribeChannel, telegramClient, update.Message)
	case commands.Subscribe.Matches(messageText):
		commands.Handle(commands.Subscribe, telegramClient, update.Message)
	case commands.Unsubscribe.Matches(messageText):
		commands.Handle(commands.Unsubscribe, telegramClient, update.Message)
	case commands.SetStopWordsChannel.Matches(messageText):
		commands.Handle(commands.SetStopWordsChannel, telegramClient, update.Message)
	case commands.SetStopWords.Matches(messageText):
		commands.Handle(commands.SetStopWords, telegramClient, update.Message)
	case commands.Usage.Matches(messageText):
		commands.Handle(commands.Usage, telegramClient, update.Message)
	}
}
