package bot

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"regexp"
)

type ArgumentsParser func(command *Command, messageText string) ([]string, bool)

type CommandHandler func(
	command *Command,
	bot *Bot,
	args []string,
	message *tgbotapi.Message,
) error

type Command struct {
	Regexp         *regexp.Regexp
	SuccessMessage string
	UsageMessage   string
	ParseArguments ArgumentsParser
	HandleCommand  CommandHandler
}

func (cmd *Command) Matches(text string) bool {
	return cmd.Regexp.MatchString(text)
}

func (cmd *Command) Handle(bot *Bot, message *tgbotapi.Message) {
	var args []string
	var ok bool

	if cmd.ParseArguments != nil {
		args, ok = cmd.ParseArguments(cmd, message.Text)
		if !ok {
			bot.sendMessage(message.Chat.ID, cmd.UsageMessage)
			return
		}
	}
	err := cmd.HandleCommand(cmd, bot, args, message)
	if err != nil {
		bot.sendReplyMessage(message.Chat.ID, message.MessageID, err.Error())
	}
}

const (
	UnauthorizedError = "You have to be channel administrator to do that."
)

var (
	SpaceRegexp = regexp.MustCompile(`\s`)

	SubscribeCommand = Command{
		Regexp:         regexp.MustCompile(`\s*/2ch_subscribe\s+([\w\s]*)`),
		SuccessMessage: "Successfully subscribed!",
		UsageMessage:   "/2ch_subscribe <board1> <board2>...",
		ParseArguments: func(cmd *Command, messageText string) ([]string, bool) {
			args := cmd.Regexp.FindStringSubmatch(messageText)
			return args, len(args) > 1 && args[1] != ""
		},
		HandleCommand: func(
			cmd *Command,
			bot *Bot,
			args []string,
			message *tgbotapi.Message,
		) error {
			chatID := message.Chat.ID
			if !bot.isUserAdministrator(chatID, message.From.ID) {
				return errors.New(UnauthorizedError)
			}
			boardNames := SpaceRegexp.Split(args[1], -1)
			subscribeToBoards(bot, chatID, chatID, boardNames, cmd.SuccessMessage)
			log.Printf("[%s] subscribed to %v", message.From.UserName, boardNames)
			return nil
		},
	}
	UnsubscribeCommand = Command{
		Regexp:         regexp.MustCompile(`\s*/2ch_unsubscribe\s*([\w\s]*)?`),
		SuccessMessage: "Successfully unsubscribed!",
		UsageMessage:   "/2ch_unsubscribe <board1> <board2>...",
		ParseArguments: func(cmd *Command, messageText string) ([]string, bool) {
			args := cmd.Regexp.FindStringSubmatch(messageText)
			return args, true
		},
		HandleCommand: func(
			cmd *Command,
			bot *Bot,
			args []string,
			message *tgbotapi.Message,
		) error {
			chatID := message.Chat.ID
			if !bot.isUserAdministrator(chatID, message.From.ID) {
				return errors.New(UnauthorizedError)
			}
			boardNames := SpaceRegexp.Split(args[1], -1)
			unsubscribeFromBoards(bot, chatID, chatID, boardNames, cmd.SuccessMessage)
			log.Printf("%v [%s] unsubscribed from %v", chatID, message.From.UserName, boardNames)
			return nil
		},
	}
	UsageCommand = Command{
		Regexp: regexp.MustCompile(`\s*/2ch_usage`),
		SuccessMessage: SubscribeCommand.UsageMessage +
			"\n" +
			UnsubscribeCommand.UsageMessage +
			"\n" +
			SubscribeChannelCommand.UsageMessage,
		HandleCommand: func(
			cmd *Command,
			bot *Bot,
			args []string,
			message *tgbotapi.Message,
		) error {
			bot.sendMessage(message.Chat.ID, cmd.SuccessMessage)
			return nil
		},
	}
	SubscribeChannelCommand = Command{
		Regexp:         regexp.MustCompile(`\s*/2ch_subscribe_channel\s+(@\w*)\s*([\w\s]*)`),
		SuccessMessage: "Successfully subscribed!",
		UsageMessage:   "/2ch_subscribe_channel @<channel> <board1> <board2>...",
		ParseArguments: func(cmd *Command, messageText string) ([]string, bool) {
			args := cmd.Regexp.FindStringSubmatch(messageText)
			return args, len(args) > 2 && args[1] != "" && args[2] != ""
		},
		HandleCommand: func(
			cmd *Command,
			bot *Bot,
			args []string,
			message *tgbotapi.Message,
		) error {
			channelName := args[1]
			responseChatID := message.Chat.ID
			chatID, err := bot.getChatIDByChannelName(channelName)
			if err != nil {
				bot.sendMessage(responseChatID, cmd.UsageMessage)
				return err
			}
			if !bot.isUserAdministrator(chatID, message.From.ID) {
				return errors.New(UnauthorizedError)
			}
			boardNames := SpaceRegexp.Split(args[2], -1)

			subscribeToBoards(bot, chatID, responseChatID, boardNames, cmd.SuccessMessage)
			log.Printf(
				"[%s] subscribed channel %s to %v",
				message.From.UserName,
				channelName,
				boardNames,
			)
			return nil
		},
	}
	UnsubscribeChannelCommand = Command{
		Regexp:         regexp.MustCompile(`\s*/2ch_unsubscribe_channel\s+(@\w*)\s*([\w\s]*)?`),
		SuccessMessage: "Successfully unsubscribed!",
		UsageMessage:   "/2ch_unsubscribe_channel @<channel> <board1> <board2>...",
		ParseArguments: func(cmd *Command, messageText string) ([]string, bool) {
			args := cmd.Regexp.FindStringSubmatch(messageText)
			return args, len(args) > 2 && args[1] != ""
		},
		HandleCommand: func(
			cmd *Command,
			bot *Bot,
			args []string,
			message *tgbotapi.Message,
		) error {
			channelName := args[1]
			responseChatID := message.Chat.ID
			chatID, err := bot.getChatIDByChannelName(channelName)
			if err != nil {
				bot.sendMessage(responseChatID, cmd.UsageMessage)
				return err
			}
			if !bot.isUserAdministrator(chatID, message.From.ID) {
				return errors.New(UnauthorizedError)
			}
			boardNames := SpaceRegexp.Split(args[2], -1)

			unsubscribeFromBoards(bot, chatID, responseChatID, boardNames, cmd.SuccessMessage)
			log.Printf(
				"[%s] unsubscribed channel %s from %v",
				message.From.UserName,
				channelName,
				boardNames,
			)
			return nil
		},
	}
)

func subscribeToBoards(bot *Bot,
	subscribeChatID int64,
	responseChatID int64,
	boardNames []string,
	successMessage string) {
	for _, boardName := range boardNames {
		if boardName == "" {
			continue
		}
		bot.storage.SubscribeChat(boardName, subscribeChatID)
	}
	bot.sendMessage(responseChatID, successMessage)
}

func unsubscribeFromBoards(bot *Bot,
	unsubscribeChatID int64,
	responseChatID int64,
	boardNames []string,
	successMessage string) {
	for _, boardName := range boardNames {
		bot.storage.UnsubscribeChat(boardName, unsubscribeChatID)
	}
	bot.sendMessage(responseChatID, successMessage)
}

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
	case SubscribeChannelCommand.Matches(messageText):
		SubscribeChannelCommand.Handle(bot, update.Message)
	case UnsubscribeChannelCommand.Matches(messageText):
		UnsubscribeChannelCommand.Handle(bot, update.Message)
	case SubscribeCommand.Matches(messageText):
		SubscribeCommand.Handle(bot, update.Message)
	case UnsubscribeCommand.Matches(messageText):
		UnsubscribeCommand.Handle(bot, update.Message)
	case UsageCommand.Matches(messageText):
		UsageCommand.Handle(bot, update.Message)
	}
}

func (bot *Bot) getChatIDByChannelName(channelName string) (int64, error) {
	chatConfig := tgbotapi.ChatConfig{SuperGroupUsername: channelName}
	chat, err := bot.telegramApi.GetChat(chatConfig)
	if err != nil {
		return 0, err
	}
	return chat.ID, nil
}

func (bot *Bot) isUserAdministrator(chatID int64, userID int) bool {
	chatConfig := tgbotapi.ChatConfig{ChatID: chatID}
	chat, err := bot.telegramApi.GetChat(chatConfig)
	switch {
	case err != nil:
		return false
	case chat.IsPrivate():
		return true
	}

	members, err := bot.telegramApi.GetChatAdministrators(chatConfig)
	if err != nil {
		return false
	}
	for _, member := range members {
		if member.User.ID == userID {
			return true
		}
	}
	return false
}

func (bot *Bot) sendReplyMessage(chatID int64, messageID int, messageText string) {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyToMessageID = messageID

	bot.telegramApi.Send(msg)
}

func (bot *Bot) sendMessage(chatID int64, messageText string) {
	msg := tgbotapi.NewMessage(chatID, messageText)

	bot.telegramApi.Send(msg)
}
