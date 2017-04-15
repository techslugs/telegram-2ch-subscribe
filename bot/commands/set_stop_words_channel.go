package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
	"log"
	"fmt"
	"regexp"
)

func BuildSetStopWordsChannel(botName string) Command {
	regexp_template := `(?s)\s*/set_stop_words_channel(?:@%s)?\s+(@\w*)\s+(\w+)(.*)`
	regexp_source := fmt.Sprintf(regexp_template, botName)
	return &SetStopWordsChannelCommand{
		BaseCommand{
			regexp:         regexp.MustCompile(regexp_source),
			successMessage: "Successfully set stop words!",
			usageMessage:   "/set_stop_words_channel @channel_name <board>\n\tstop words\n\t1 per line",
		},
	}
}

type SetStopWordsChannelCommand struct {
	BaseCommand
}

func (cmd *SetStopWordsChannelCommand) Parse(messageText string) ([]string, bool) {
	args := cmd.regexp.FindStringSubmatch(messageText)
	return args, len(args) > 2 && args[1] != "" && args[2] != ""
}

func (cmd *SetStopWordsChannelCommand) Process(
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
	stopWords := parseStopWords(args[3])

	boardNames, err :=
		telegramClient.SetStopWords(
			chatID,
			responseChatID,
			args[2],
			stopWords,
			cmd.SuccessMessage(),
		)

	if err != nil {
		log.Printf(
			"[%s] Error %s while setting stop words to %v",
			message.From.UserName,
			err,
			boardNames,
		)
		return errors.New(UnknownError)
	}
	log.Printf(
		"[%s] set stop words for channel %s, %v",
		message.From.UserName,
		channelName,
		boardNames,
	)
	return nil
}
