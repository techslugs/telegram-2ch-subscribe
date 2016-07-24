package handlers

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
)

var SubscribeChannelCommand = Command{
	Regexp:         regexp.MustCompile(`\s*/2ch_subscribe_channel\s+(@\w*)\s*([\w\s]*)`),
	SuccessMessage: "Successfully subscribed!",
	UsageMessage:   "/2ch_subscribe_channel @<channel> <board1> <board2>...",
	ParseArguments: func(cmd *Command, messageText string) ([]string, bool) {
		args := cmd.Regexp.FindStringSubmatch(messageText)
		return args, len(args) > 2 && args[1] != "" && args[2] != ""
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
			return errors.New(cmd.UsageMessage)
		}
		if !telegramCommands.IsUserAdministrator(chatID, message.From.ID) {
			return errors.New(UnauthorizedError)
		}

		boardNames, err := telegramCommands.
			SubscribeToBoards(chatID, responseChatID, args[2], cmd.SuccessMessage)

		if err != nil {
			log.Printf(
				"[%s] Error %s while subscribing %s to %v",
				message.From.UserName,
				err.Error(),
				channelName,
				boardNames,
			)
			return errors.New(UnknownError)
		}
		log.Printf(
			"[%s] subscribed channel %s to %v",
			message.From.UserName,
			channelName,
			boardNames,
		)
		return nil
	},
}
