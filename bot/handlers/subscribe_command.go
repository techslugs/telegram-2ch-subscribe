package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
)

var SubscribeCommand = Command{
	Regexp:         regexp.MustCompile(`\s*/2ch_subscribe\s+([\w\s]*)`),
	SuccessMessage: "Successfully subscribed!",
	UsageMessage:   "/2ch_subscribe <board1> <board2>...",
	ParseArguments: parseSubscribeCommandArguments,
	HandleCommand:  handleSubscribeCommand,
}

func parseSubscribeCommandArguments(cmd *Command, messageText string) ([]string, bool) {
	args := cmd.Regexp.FindStringSubmatch(messageText)
	return args, len(args) > 1 && args[1] != ""
}

func handleSubscribeCommand(
	cmd *Command,
	telegramCommands *TelegramCommandsHandler,
	args []string,
	message *tgbotapi.Message,
) error {
	chatID := message.Chat.ID
	if !telegramCommands.IsUserAdministrator(chatID, message.From.ID) {
		return errors.New(UnauthorizedError)
	}

	boardNames, err := telegramCommands.
		SubscribeToBoards(chatID, chatID, args[1], cmd.SuccessMessage)

	if err != nil {
		log.Printf(
			"[%s] Error %s while subscribing to %v",
			message.From.UserName,
			err,
			boardNames,
		)
		return errors.New(UnknownError)
	}
	log.Printf("[%s] subscribed to %v", message.From.UserName, boardNames)
	return nil
}
