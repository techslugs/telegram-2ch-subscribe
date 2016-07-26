package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tmwh/telegram-2ch-subscribe/storage"
	"log"
	"regexp"
)

var (
	SpaceRegexp = regexp.MustCompile(`\s`)
)

type Client struct {
	TelegramAPI     *tgbotapi.BotAPI
	TelegramUpdates <-chan tgbotapi.Update
	Storage         *storage.Storage
}

func NewClient(telegramToken string, storage *storage.Storage) (*Client, error) {
	api, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	telegramUpdates, err := api.GetUpdatesChan(updateConfig)
	if err != nil {
		return nil, err
	}

	client := &Client{
		TelegramAPI:     api,
		TelegramUpdates: telegramUpdates,
		Storage:         storage,
	}
	return client, nil
}

func (client *Client) GetChatIDByChannelName(
	channelName string,
) (int64, error) {
	chatConfig := tgbotapi.ChatConfig{SuperGroupUsername: channelName}
	chat, err := client.TelegramAPI.GetChat(chatConfig)
	if err != nil {
		return 0, err
	}
	return chat.ID, nil
}

func (client *Client) IsUserAdministrator(
	chatID int64,
	userID int,
) bool {
	chatConfig := tgbotapi.ChatConfig{ChatID: chatID}
	chat, err := client.TelegramAPI.GetChat(chatConfig)
	switch {
	case err != nil:
		return false
	case chat.IsPrivate():
		return true
	}

	members, err := client.TelegramAPI.GetChatAdministrators(chatConfig)
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

func (client *Client) SendReplyMessage(
	chatID int64,
	messageID int,
	messageText string,
) {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyToMessageID = messageID

	_, err := client.TelegramAPI.Send(msg)
	if err != nil {
		client.reportChatError(chatID, err)
	}
}

func (client *Client) SendMessage(
	chatID int64,
	messageText string,
) {
	msg := tgbotapi.NewMessage(chatID, messageText)

	_, err := client.TelegramAPI.Send(msg)
	if err != nil {
		client.reportChatError(chatID, err)
	} else {
		client.reportChatSuccess(chatID)
	}
}

func (client *Client) SendMarkdownMessage(
	chatID int64,
	messageText string,
) {
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "Markdown"

	_, err := client.TelegramAPI.Send(msg)
	if err != nil {
		client.reportChatError(chatID, err)
	} else {
		client.reportChatSuccess(chatID)
	}
}

func (client *Client) reportChatError(chatID int64, err error) {
	if err := client.Storage.ReportChatError(chatID); err != nil {
		log.Printf("Error while reporting Chat [%v] error: %s", chatID, err)
		return
	}
	log.Printf("Chat [%v] error: %s", chatID, err)
}

func (client *Client) reportChatSuccess(chatID int64) {
	if err := client.Storage.ReportChatSuccess(chatID); err != nil {
		log.Printf("Error while reporting Chat [%v] error: %s", chatID, err)
		return
	}
}

func (client *Client) SubscribeToBoards(
	subscribeChatID int64,
	responseChatID int64,
	boardNamesSplitBySpace string,
	successMessage string,
) ([]string, error) {
	boardNames, err := client.Storage.
		FilterInvalidBoardNames(SpaceRegexp.Split(boardNamesSplitBySpace, -1))
	if err != nil {
		return boardNames, err
	}

	for _, boardName := range boardNames {
		if boardName == "" {
			continue
		}
		client.Storage.SubscribeChat(boardName, subscribeChatID)
	}
	client.SendMessage(responseChatID, successMessage)
	return boardNames, nil
}

func (client *Client) UnsubscribeFromBoards(
	unsubscribeChatID int64,
	responseChatID int64,
	boardNames []string,
	successMessage string,
) {
	for _, boardName := range boardNames {
		client.Storage.UnsubscribeChat(boardName, unsubscribeChatID)
	}
	client.SendMessage(responseChatID, successMessage)
}
