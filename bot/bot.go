package bot

import (
	"log"
	"net/http"
	"github.com/tmwh/telegram-2ch-subscribe/bot/fetchers"
	"github.com/tmwh/telegram-2ch-subscribe/storage"
	"github.com/tmwh/telegram-2ch-subscribe/telegram"
	"time"
)

type Bot struct {
	Storage           *storage.Storage
	TelegramClient    *telegram.Client
	BoardListFetcher  *fetchers.BoardsListFetcher
	BoardsInfoFetcher *fetchers.BoardsInfoFetcher
}

func StartBot(
	boardsListUpdateTimeout time.Duration,
	boardInfoUpdateTimeout time.Duration,
	telegramClient *telegram.Client,
	storage *storage.Storage,
) error {

	bot := Bot{Storage: storage, TelegramClient: telegramClient}

	httpClient := &http.Client{Timeout: time.Second * 5}
	bot.BoardListFetcher = &fetchers.BoardsListFetcher{HttpClient: httpClient}
	bot.BoardsInfoFetcher = &fetchers.BoardsInfoFetcher{
		HttpClient: httpClient,
		Storage:    bot.Storage,
	}

	go StartHandleCommandsFromTelegram(bot.TelegramClient)

	go fetchers.StartFetchingBoardsListUpdates(
		bot.BoardListFetcher,
		boardsListUpdateTimeout,
		func(boardInfo *fetchers.BoardsList, err error) {
			if err != nil {
				log.Printf(`Error: Could not fetch boards list. %s`, err)
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
				log.Printf(`Error: Could not fetch board info for "%s". %s`, boardInfo.Board, err)
				return
			}

			bot.publishBoardInfo(boardInfo)
		})

	return nil
}
