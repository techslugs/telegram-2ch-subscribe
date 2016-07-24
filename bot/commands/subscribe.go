package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
	"github.com/tmwh/telegram-2ch-subscribe/telegram"
)

var Subscribe = &SubscribeCommand{
	BaseCommand{
		regexp:         regexp.MustCompile(`\s*/2ch_subscribe\s+([\w\s]*)`),
		successMessage: "Successfully subscribed!",
		usageMessage:   "/2ch_subscribe <board1> <board2>...",
	},
}

type SubscribeCommand struct {
	BaseCommand
}

func (cmd *SubscribeCommand) Parse(messageText string) ([]string, bool) {
	args := cmd.regexp.FindStringSubmatch(messageText)
	return args, len(args) > 1 && args[1] != ""
}

func (cmd *SubscribeCommand) Process(
	telegramClient *telegram.Client,
	args []string,
	message *tgbotapi.Message,
) error {
	chatID := message.Chat.ID
	if !telegramClient.IsUserAdministrator(chatID, message.From.ID) {
		return errors.New(UnauthorizedError)
	}

	boardNames, err :=
		telegramClient.SubscribeToBoards(chatID, chatID, args[1], cmd.SuccessMessage())

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
