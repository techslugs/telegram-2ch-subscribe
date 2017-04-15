package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"fmt"
	"regexp"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
)

func BuildUnsubscribe(botName string) Command {
	regexp_template := `\s*/unsubscribe(?:@%s)?\s*([\w\s]*)?`
	regexp_source := fmt.Sprintf(regexp_template, botName)
  return &UnsubscribeCommand{
  	BaseCommand{
  		regexp:         regexp.MustCompile(regexp_source),
  		successMessage: "Successfully unsubscribed!",
  		usageMessage:   "/unsubscribe <board1> <board2>...",
  	},
  }
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
