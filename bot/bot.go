package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"time"
)

type Bot struct {
	telegramApi     *tgbotapi.BotAPI
	telegramUpdates <-chan tgbotapi.Update
	httpClient      *http.Client
}

func StartBot(telegramToken string, updatesTimeout time.Duration) error {
	bot := Bot{}
	err := bot.setupTelegramApi(telegramToken)
	if err != nil {
		return err
	}

	bot.httpClient = &http.Client{
		Timeout: time.Second * 5,
	}

	go bot.handleCommandsFromTelegram()
	go bot.fetchBoardInfoUpdates(updatesTimeout, func(boardInfo *BoardInfo, err error) {
		if err != nil {
			return
		}

		bot.publishBoardInfo(boardInfo)
	})

	return nil
}

func (bot *Bot) setupTelegramApi(telegramToken string) error {
	api, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	telegramUpdates, err := api.GetUpdatesChan(updateConfig)
	if err != nil {
		return err
	}

	bot.telegramApi = api
	bot.telegramUpdates = telegramUpdates
	return nil
}
