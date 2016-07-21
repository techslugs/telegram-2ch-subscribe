package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
)

type Command struct {
	Regexp         *regexp.Regexp
	SuccessMessage string
	UsageMessage   string
}

func (cmd *Command) Matches(text string) bool {
	return cmd.Regexp.MatchString(text)
}

var (
	SpaceRegexp = regexp.MustCompile(`\s`)

	SubscribeCommand = Command{
		Regexp:         regexp.MustCompile(`\s*/2ch_subscribe\s?([\w\s]*)`),
		SuccessMessage: "Successfully subscribed!",
		UsageMessage:   "/2ch_subscribe <board1> <board2>...",
	}
	UnsubscribeCommand = Command{
		Regexp:         regexp.MustCompile(`\s*/2ch_unsubscribe\s?([\w\s]*)?`),
		SuccessMessage: "Successfully unsubscribed!",
		UsageMessage:   "/2ch_unsubscribe <board1> <board2>...",
	}
	UsageCommand = Command{
		Regexp:         regexp.MustCompile(`\s*/2ch_usage`),
		SuccessMessage: SubscribeCommand.UsageMessage + "\n" + UnsubscribeCommand.UsageMessage,
	}
	SubscribeChannelRegexp = regexp.MustCompile(
		`\s*/2ch_subscribe-channel\s?([\w\s]*)`,
	)
)

func (bot *Bot) handleCommandsFromTelegram() {
	for update := range bot.telegramUpdates {
		bot.parseAndHandleCommand(&update)
	}
}

func (bot *Bot) parseAndHandleCommand(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	messageText := update.Message.Text
	switch {
	case SubscribeCommand.Matches(messageText):
		bot.parseAndHandleSubscribe(update.Message)
	case UnsubscribeCommand.Matches(messageText):
		bot.parseAndHandleUnsubscribe(update.Message)
	case UsageCommand.Matches(messageText):
		bot.parseAndHandleUsage(update.Message)
	}
}

func (bot *Bot) parseAndHandleSubscribe(message *tgbotapi.Message) {
	matches := SubscribeCommand.Regexp.FindStringSubmatch(message.Text)
	if len(matches) < 2 || matches[1] == "" {
		bot.sendMessage(message.Chat.ID, SubscribeCommand.UsageMessage)
		return
	}
	chatID := message.Chat.ID
	boardNames := SpaceRegexp.Split(matches[1], -1)
	for _, boardName := range boardNames {
		if boardName == "" {
			continue
		}
		bot.storage.SubscribeChat(boardName, chatID)
	}
	log.Printf("[%s] subscribed to %v", message.From.UserName, boardNames)
	bot.sendMessage(chatID, SubscribeCommand.SuccessMessage)
}

func (bot *Bot) parseAndHandleUnsubscribe(message *tgbotapi.Message) {
	matches := UnsubscribeCommand.Regexp.FindStringSubmatch(message.Text)
	var boardNames []string
	if len(matches) < 2 {
		boardNames = []string{""}
	} else {
		boardNames = SpaceRegexp.Split(matches[1], -1)
	}

	chatID := message.Chat.ID
	for _, boardName := range boardNames {
		bot.storage.UnsubscribeChat(boardName, chatID)
	}
	log.Printf("%v [%s] unsubscribed from %v", chatID, message.From.UserName, boardNames)
	bot.sendMessage(chatID, UnsubscribeCommand.SuccessMessage)
}

func (bot *Bot) parseAndHandleUsage(message *tgbotapi.Message) {
	bot.sendMessage(message.Chat.ID, UsageCommand.SuccessMessage)
}

func (bot *Bot) sendMessage(chatID int64, messageText string) {
	msg := tgbotapi.NewMessage(chatID, messageText)

	bot.telegramApi.Send(msg)
}
