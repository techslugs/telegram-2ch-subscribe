package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"telegram-2ch-news-bot/bot/fetchers"
	"telegram-2ch-news-bot/bot/handlers"
	"telegram-2ch-news-bot/storage"
	"time"
)

type Bot struct {
	Storage           *storage.Storage
	TelegramCommands  *handlers.TelegramCommandsHandler
	BoardListFetcher  *fetchers.BoardsListFetcher
	BoardsInfoFetcher *fetchers.BoardsInfoFetcher
}

func StartBot(
	telegramToken string,
	boardsListUpdateTimeout time.Duration,
	boardInfoUpdateTimeout time.Duration,
	storage *storage.Storage,
) error {

	bot := Bot{Storage: storage}
	err := bot.setupTelegramCommands(telegramToken)
	if err != nil {
		return err
	}

	httpClient := &http.Client{Timeout: time.Second * 5}
	bot.BoardListFetcher = &fetchers.BoardsListFetcher{HttpClient: httpClient}
	bot.BoardsInfoFetcher = &fetchers.BoardsInfoFetcher{
		HttpClient: httpClient,
		Storage:    bot.Storage,
	}

	go handlers.StartHandleCommandsFromTelegram(bot.TelegramCommands)
	go fetchers.StartFetchingBoardsListUpdates(
		bot.BoardListFetcher,
		boardsListUpdateTimeout,
		func(boardInfo *fetchers.BoardsList, err error) {
			if err != nil {
				log.Printf(
					`Error: Could not fetch boards list. %s`,
					err.Error())
				return
			}

			bot.saveBoardsList(boardInfo)
			log.Println("Updated boards list.")
		})

	go fetchers.StartFetchingBoardInfoUpdates(
		bot.BoardsInfoFetcher,
		boardInfoUpdateTimeout,
		func(boardInfo *fetchers.BoardInfo, err error) {
			if err != nil {
				log.Printf(
					`Error: Could not fetch board info for "%s". %s`,
					boardInfo.Board,
					err.Error())
				return
			}

			bot.publishBoardInfo(boardInfo)
		})

	return nil
}

func (bot *Bot) setupTelegramCommands(telegramToken string) error {
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

	bot.TelegramCommands = &handlers.TelegramCommandsHandler{
		TelegramAPI:     api,
		TelegramUpdates: telegramUpdates,
		Storage:         bot.Storage,
	}
	return nil
}
