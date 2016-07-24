package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
)

var UnsubscribeCommand = Command{
	Regexp:         regexp.MustCompile(`\s*/2ch_unsubscribe\s*([\w\s]*)?`),
	SuccessMessage: "Successfully unsubscribed!",
	UsageMessage:   "/2ch_unsubscribe <board1> <board2>...",
	ParseArguments: parseUnsubscribeCommandArguments,
	HandleCommand:  handleUnsubscribeCommand,
}

func parseUnsubscribeCommandArguments(cmd *Command, messageText string) ([]string, bool) {
	args := cmd.Regexp.FindStringSubmatch(messageText)
	return args, true
}

func handleUnsubscribeCommand(
	cmd *Command,
	telegramCommands *TelegramCommandsHandler,
	args []string,
	message *tgbotapi.Message,
) error {
	chatID := message.Chat.ID
	if !telegramCommands.IsUserAdministrator(chatID, message.From.ID) {
		return errors.New(UnauthorizedError)
	}
	boardNames := SpaceRegexp.Split(args[1], -1)
	telegramCommands.UnsubscribeFromBoards(chatID, chatID, boardNames, cmd.SuccessMessage)
	log.Printf("%v [%s] unsubscribed from %v", chatID, message.From.UserName, boardNames)
	return nil
}
