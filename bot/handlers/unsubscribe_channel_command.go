package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
)

var UnsubscribeChannelCommand = Command{
	Regexp:         regexp.MustCompile(`\s*/2ch_unsubscribe_channel\s+(@\w*)\s*([\w\s]*)?`),
	SuccessMessage: "Successfully unsubscribed!",
	UsageMessage:   "/2ch_unsubscribe_channel @<channel> <board1> <board2>...",
	ParseArguments: func(cmd *Command, messageText string) ([]string, bool) {
		args := cmd.Regexp.FindStringSubmatch(messageText)
		return args, len(args) > 2 && args[1] != ""
	},
	HandleCommand: func(
		cmd *Command,
		telegramCommands *TelegramCommandsHandler,
		args []string,
		message *tgbotapi.Message,
	) error {
		channelName := args[1]
		responseChatID := message.Chat.ID
		chatID, err := telegramCommands.GetChatIDByChannelName(channelName)
		if err != nil {
			telegramCommands.SendMessage(responseChatID, cmd.UsageMessage)
			return err
		}
		if !telegramCommands.IsUserAdministrator(chatID, message.From.ID) {
			return errors.New(UnauthorizedError)
		}
		boardNames := SpaceRegexp.Split(args[2], -1)

		telegramCommands.UnsubscribeFromBoards(chatID, responseChatID, boardNames, cmd.SuccessMessage)
		log.Printf(
			"[%s] unsubscribed channel %s from %v",
			message.From.UserName,
			channelName,
			boardNames,
		)
		return nil
	},
}
