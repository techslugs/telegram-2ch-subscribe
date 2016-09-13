package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
	"log"
	"regexp"
	"strconv"
)

var SubscribeChannel = &SubscribeChannelCommand{
	BaseCommand{
		regexp:         regexp.MustCompile(`\s*/2ch_subscribe_channel\s+(@\w*)\s*(\w+)\s*(\d+(\.\d*)?)`),
		successMessage: "Successfully subscribed!",
		usageMessage:   "/2ch_subscribe_channel @<channel> <board> <min score>",
	},
}

type SubscribeChannelCommand struct {
	BaseCommand
}

func (cmd *SubscribeChannelCommand) Parse(messageText string) ([]string, bool) {
	args := cmd.regexp.FindStringSubmatch(messageText)
	return args, len(args) > 2 && args[1] != "" && args[2] != "" && args[3] != ""
}

func (cmd *SubscribeChannelCommand) Process(
	telegramClient *telegram.Client,
	args []string,
	message *tgbotapi.Message,
) error {
	channelName := args[1]
	responseChatID := message.Chat.ID
	chatID, err := telegramClient.GetChatIDByChannelName(channelName)
	if err != nil {
		return errors.New(cmd.UsageMessage())
	}
	if !telegramClient.IsUserAdministrator(chatID, message.From.ID) {
		return errors.New(UnauthorizedError)
	}

	minScore, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return errors.New(cmd.UsageMessage())
	}

	boardNames, err := telegramClient.
		SubscribeToBoards(chatID, responseChatID, args[2], minScore, cmd.SuccessMessage())

	if err != nil {
		log.Printf(
			"[%s] Error %s while subscribing %s to %v",
			message.From.UserName,
			err,
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
}
