package handlers

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"telegram-2ch-news-bot/storage"
)

type TelegramCommandsHandler struct {
	TelegramAPI     *tgbotapi.BotAPI
	TelegramUpdates <-chan tgbotapi.Update
	Storage         *storage.Storage
}

func StartHandleCommandsFromTelegram(handler *TelegramCommandsHandler) {
	for update := range handler.TelegramUpdates {
		handler.parseAndHandleCommand(&update)
	}
}

func (handler *TelegramCommandsHandler) parseAndHandleCommand(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	messageText := update.Message.Text
	switch {
	case SubscribeChannelCommand.Matches(messageText):
		SubscribeChannelCommand.Handle(handler, update.Message)
	case UnsubscribeChannelCommand.Matches(messageText):
		UnsubscribeChannelCommand.Handle(handler, update.Message)
	case SubscribeCommand.Matches(messageText):
		SubscribeCommand.Handle(handler, update.Message)
	case UnsubscribeCommand.Matches(messageText):
		UnsubscribeCommand.Handle(handler, update.Message)
	case UsageCommand.Matches(messageText):
		UsageCommand.Handle(handler, update.Message)
	}
}

func (handler *TelegramCommandsHandler) GetChatIDByChannelName(
	channelName string,
) (int64, error) {
	chatConfig := tgbotapi.ChatConfig{SuperGroupUsername: channelName}
	chat, err := handler.TelegramAPI.GetChat(chatConfig)
	if err != nil {
		return 0, err
	}
	return chat.ID, nil
}

func (handler *TelegramCommandsHandler) IsUserAdministrator(
	chatID int64,
	userID int,
) bool {
	chatConfig := tgbotapi.ChatConfig{ChatID: chatID}
	chat, err := handler.TelegramAPI.GetChat(chatConfig)
	switch {
	case err != nil:
		return false
	case chat.IsPrivate():
		return true
	}

	members, err := handler.TelegramAPI.GetChatAdministrators(chatConfig)
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

func (handler *TelegramCommandsHandler) SendReplyMessage(
	chatID int64,
	messageID int,
	messageText string,
) {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyToMessageID = messageID

	_, err := handler.TelegramAPI.Send(msg)
	if err != nil {
		handler.reportChatError(chatID, err)
	}
}

func (handler *TelegramCommandsHandler) SendMessage(
	chatID int64,
	messageText string,
) {
	msg := tgbotapi.NewMessage(chatID, messageText)

	_, err := handler.TelegramAPI.Send(msg)
	if err != nil {
		handler.reportChatError(chatID, err)
	}
}

func (handler *TelegramCommandsHandler) reportChatError(chatID int64, err error) {
	if err := handler.Storage.ReportChatError(chatID); err != nil {
		log.Printf("Error while reporting Chat [%v] error: %s", chatID, err)
		return
	}
	log.Printf("Chat [%v] error: %s", chatID, err)
}

func (handler *TelegramCommandsHandler) SubscribeToBoards(
	subscribeChatID int64,
	responseChatID int64,
	boardNamesSplitBySpace string,
	successMessage string,
) ([]string, error) {
	boardNames, err := handler.Storage.
		FilterInvalidBoardNames(SpaceRegexp.Split(boardNamesSplitBySpace, -1))
	if err != nil {
		return boardNames, err
	}

	for _, boardName := range boardNames {
		if boardName == "" {
			continue
		}
		handler.Storage.SubscribeChat(boardName, subscribeChatID)
	}
	handler.SendMessage(responseChatID, successMessage)
	return boardNames, nil
}

func (handler *TelegramCommandsHandler) UnsubscribeFromBoards(
	unsubscribeChatID int64,
	responseChatID int64,
	boardNames []string,
	successMessage string,
) {
	for _, boardName := range boardNames {
		handler.Storage.UnsubscribeChat(boardName, unsubscribeChatID)
	}
	handler.SendMessage(responseChatID, successMessage)
}
