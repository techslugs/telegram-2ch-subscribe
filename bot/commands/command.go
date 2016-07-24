package commands

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
	"telegram-2ch-news-bot/telegram"
)

const (
	UnauthorizedError = "You have to be channel administrator to do that."
	UnknownError      = "Something went wrong, we are trying to fix it ASAP."
)

var (
	SpaceRegexp = regexp.MustCompile(`\s`)
)

type Command interface {
	UsageMessage() string
	SuccessMessage() string
	Matches(text string) bool
	Parse(messageText string) ([]string, bool)
	Process(telegramClient *telegram.Client, args []string, message *tgbotapi.Message) error
}

func Handle(cmd Command, telegramClient *telegram.Client, message *tgbotapi.Message) {
	var args []string
	var ok bool

	args, ok = cmd.Parse(message.Text)
	if !ok {
		telegramClient.SendMessage(message.Chat.ID, cmd.UsageMessage())
		return
	}
	err := cmd.Process(telegramClient, args, message)
	if err != nil {
		telegramClient.SendReplyMessage(message.Chat.ID, message.MessageID, err.Error())
	}
}
