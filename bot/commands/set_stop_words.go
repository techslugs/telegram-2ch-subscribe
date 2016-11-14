package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
	"log"
	"regexp"
	"strings"
)

var SetStopWords = &SetStopWordsCommand{
	BaseCommand{
		regexp:         regexp.MustCompile(`(?s)\s*/2ch_set_stop_words\s+(\w+)(.*)`),
		successMessage: "Successfully set stop words!",
		usageMessage:   "/2ch_set_stop_words <board>\n\tstop words\n\t1 per line",
	},
}

type SetStopWordsCommand struct {
	BaseCommand
}

func (cmd *SetStopWordsCommand) Parse(messageText string) ([]string, bool) {
	args := cmd.regexp.FindStringSubmatch(messageText)
	return args, len(args) > 1 && args[1] != ""
}

func (cmd *SetStopWordsCommand) Process(
	telegramClient *telegram.Client,
	args []string,
	message *tgbotapi.Message,
) error {
	chatID := message.Chat.ID
	if !telegramClient.IsUserAdministrator(chatID, message.From.ID) {
		return errors.New(UnauthorizedError)
	}
	stopWords := parseStopWords(args[2])

	boardNames, err :=
		telegramClient.SetStopWords(
			chatID,
			chatID,
			args[1],
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
	log.Printf("[%s] set stop words for %v", message.From.UserName, boardNames)
	return nil
}

func parseStopWords(message string) []string {
	words := strings.Split(message, "\n")
	stopWords := words[:0]
	for _, word := range words {
		if word != "" {
			stopWords = append(stopWords, word)
		}
	}

	return stopWords
}
