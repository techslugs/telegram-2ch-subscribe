package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
	"telegram-2ch-news-bot/telegram"
)

var Unsubscribe = &UnsubscribeCommand{
	BaseCommand{
		regexp:         regexp.MustCompile(`\s*/2ch_unsubscribe\s*([\w\s]*)?`),
		successMessage: "Successfully unsubscribed!",
		usageMessage:   "/2ch_unsubscribe <board1> <board2>...",
	},
}

type UnsubscribeCommand struct {
	BaseCommand
}

func (cmd *UnsubscribeCommand) Parse(messageText string) ([]string, bool) {
	args := cmd.regexp.FindStringSubmatch(messageText)
	return args, true
}

func (cmd *UnsubscribeCommand) Process(
	telegramClient *telegram.Client,
	args []string,
	message *tgbotapi.Message,
) error {
	chatID := message.Chat.ID
	if !telegramClient.IsUserAdministrator(chatID, message.From.ID) {
		return errors.New(UnauthorizedError)
	}
	boardNames := SpaceRegexp.Split(args[1], -1)
	telegramClient.UnsubscribeFromBoards(chatID, chatID, boardNames, cmd.SuccessMessage())
	log.Printf("%v [%s] unsubscribed from %v", chatID, message.From.UserName, boardNames)
	return nil
}
