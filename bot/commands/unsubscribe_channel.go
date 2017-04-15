package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
	"fmt"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
)

func BuildUnsubscribeChannel(botName string) Command {
	regexp_template := `\s*/unsubscribe_channel(?:@%s)?\s+(@\w*)\s*([\w\s]*)?`
	regexp_source := fmt.Sprintf(regexp_template, botName)
  return &UnsubscribeChannelCommand{
		BaseCommand{
			regexp:         regexp.MustCompile(regexp_source),
			successMessage: "Successfully unsubscribed!",
			usageMessage:   "/unsubscribe_channel @<channel> <board1> <board2>...",
		},
	}
}

type UnsubscribeChannelCommand struct {
	BaseCommand
}

func (cmd *UnsubscribeChannelCommand) Parse(messageText string) ([]string, bool) {
	args := cmd.regexp.FindStringSubmatch(messageText)
	return args, len(args) > 2 && args[1] != ""
}

func (cmd *UnsubscribeChannelCommand) Process(
	telegramClient *telegram.Client,
	args []string,
	message *tgbotapi.Message,
) error {
	channelName := args[1]
	responseChatID := message.Chat.ID
	chatID, err := telegramClient.GetChatIDByChannelName(channelName)
	if err != nil {
		telegramClient.SendMessage(responseChatID, cmd.UsageMessage())
		return err
	}
	if !telegramClient.IsUserAdministrator(chatID, message.From.ID) {
		return errors.New(UnauthorizedError)
	}
	boardNames := SpaceRegexp.Split(args[2], -1)
	telegramClient.UnsubscribeFromBoards(chatID, responseChatID, boardNames, cmd.SuccessMessage())
	log.Printf(
		"[%s] unsubscribed channel %s from %v",
		message.From.UserName,
		channelName,
		boardNames,
	)
	return nil
}
