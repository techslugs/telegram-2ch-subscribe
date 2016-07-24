package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"regexp"
)

const (
	UnauthorizedError = "You have to be channel administrator to do that."
	UnknownError      = "Something went wrong, we are trying to fix it ASAP."
)

var (
	SpaceRegexp = regexp.MustCompile(`\s`)
)

type ArgumentsParser func(command *Command, messageText string) ([]string, bool)

type CommandHandler func(
	command *Command,
	telegramCommands *TelegramCommandsHandler,
	args []string,
	message *tgbotapi.Message,
) error

type Command struct {
	Regexp         *regexp.Regexp
	SuccessMessage string
	UsageMessage   string
	ParseArguments ArgumentsParser
	HandleCommand  CommandHandler
}

func (cmd *Command) Matches(text string) bool {
	return cmd.Regexp.MatchString(text)
}

func (cmd *Command) Handle(telegramCommands *TelegramCommandsHandler, message *tgbotapi.Message) {
	var args []string
	var ok bool

	if cmd.ParseArguments != nil {
		args, ok = cmd.ParseArguments(cmd, message.Text)
		if !ok {
			telegramCommands.SendMessage(message.Chat.ID, cmd.UsageMessage)
			return
		}
	}
	err := cmd.HandleCommand(cmd, telegramCommands, args, message)
	if err != nil {
		telegramCommands.SendReplyMessage(message.Chat.ID, message.MessageID, err.Error())
	}
}
