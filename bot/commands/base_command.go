package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
)

type BaseCommand struct {
	regexp         *regexp.Regexp
	successMessage string
	usageMessage   string
}

func (cmd *BaseCommand) UsageMessage() string {
	return cmd.usageMessage
}

func (cmd *BaseCommand) SuccessMessage() string {
	return cmd.successMessage
}

func (cmd *BaseCommand) Matches(text string) bool {
	return cmd.regexp.MatchString(text)
}

func (cmd *BaseCommand) Parse(messageText string) ([]string, bool) {
	return []string{}, true
}

func (cmd *BaseCommand) Process(
	telegramClient *telegram.Client,
	args []string,
	message *tgbotapi.Message,
) error {
	telegramClient.SendMessage(message.Chat.ID, cmd.successMessage)
	return nil
}
