package commands

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/techslugs/telegram-2ch-subscribe/telegram"
	"log"
	"regexp"
	"fmt"
	"strconv"
)

func BuildSubscribe(botName string) Command {
	regexp_template := `\s*/subscribe(?:@%s)?\s+(\w+)\s*(\d+(\.\d*)?)`
	regexp_source := fmt.Sprintf(regexp_template, botName)
	return &SubscribeCommand{
		BaseCommand{
			regexp:         regexp.MustCompile(regexp_source),
			successMessage: "Successfully subscribed!",
			usageMessage:   "/subscribe <board1> <min score>",
		},
	}
}

type SubscribeCommand struct {
	BaseCommand
}

func (cmd *SubscribeCommand) Parse(messageText string) ([]string, bool) {
	args := cmd.regexp.FindStringSubmatch(messageText)
	return args, true
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
	minScore, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return errors.New(cmd.UsageMessage())
	}

	boardNames, err :=
		telegramClient.SubscribeToBoards(
			chatID,
			chatID,
			args[1],
			minScore,
			cmd.SuccessMessage(),
		)

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
