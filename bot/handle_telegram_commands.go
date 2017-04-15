package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/techslugs/telegram-2ch-subscribe/bot/commands"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
	"log"
)

func StartHandleCommandsFromTelegram(telegramClient *telegram.Client) {
	available_commands := commands.BuildCommands(telegramClient.GetMyName())
	for update := range telegramClient.TelegramUpdates {
		parseAndHandleCommand(available_commands, telegramClient, &update)
	}
}

func parseAndHandleCommand(available_commands []commands.Command, telegramClient *telegram.Client, update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	messageText := update.Message.Text
	for _, command := range available_commands {
		if command.Matches(messageText) {
			commands.Handle(command, telegramClient, update.Message)
		}
	}
}
