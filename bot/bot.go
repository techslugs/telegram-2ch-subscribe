package bot

import (
	"github.com/tmwh/telegram-2ch-subscribe/bot/fetchers"
	"github.com/tmwh/telegram-2ch-subscribe/dvach"
	"github.com/tmwh/telegram-2ch-subscribe/storage"
	"github.com/tmwh/telegram-2ch-subscribe/telegram"
	"log"
	"net/http"
	"time"
)

type Bot struct {
	Storage        *storage.Storage
	TelegramClient *telegram.Client
	DvachClient    *dvach.Client
}

func StartBot(
	boardsListUpdateTimeout time.Duration,
	boardInfoUpdateTimeout time.Duration,
	telegramClient *telegram.Client,
	storage *storage.Storage,
) error {

	bot := Bot{Storage: storage, TelegramClient: telegramClient}

	httpClient := &http.Client{Timeout: time.Second * 5}
	bot.DvachClient = &dvach.Client{HttpClient: httpClient}

	go StartHandleCommandsFromTelegram(bot.TelegramClient)

	go fetchers.StartFetchingBoardsListUpdates(
		bot.DvachClient,
		boardsListUpdateTimeout,
		func(boardInfo *dvach.BoardsList, err error) {
			if err != nil {
				log.Printf(`Error: Could not fetch boards list. %s`, err)
				return
			}

			bot.saveBoardsList(boardInfo)
			log.Println("Updated boards list.")
		})

	go fetchers.StartFetchingBoardInfoUpdates(
		bot.DvachClient,
		bot.Storage,
		boardInfoUpdateTimeout,
		func(boardInfo *dvach.BoardInfo, err error) {
			if err != nil {
				log.Printf(`Error: Could not fetch board info for "%s". %s`, boardInfo.Board, err)
				return
			}

			bot.publishBoardInfo(boardInfo)
		})

	return nil
}
