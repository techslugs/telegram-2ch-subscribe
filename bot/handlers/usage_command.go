package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
)

var UsageCommand = Command{
	Regexp: regexp.MustCompile(`\s*/2ch_usage`),
	SuccessMessage: SubscribeCommand.UsageMessage +
		"\n" +
		UnsubscribeCommand.UsageMessage +
		"\n" +
		SubscribeChannelCommand.UsageMessage,
	HandleCommand: func(
		cmd *Command,
		telegramCommands *TelegramCommandsHandler,
		args []string,
		message *tgbotapi.Message,
	) error {
		telegramCommands.SendMessage(message.Chat.ID, cmd.SuccessMessage)
		return nil
	},
}
