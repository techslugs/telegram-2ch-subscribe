package fetchers

import (
	"github.com/techslugs/telegram-2ch-subscribe/dvach"
	"github.com/techslugs/telegram-2ch-subscribe/storage"
	"log"
	"time"
)

type BoardInfoCallback func(*dvach.BoardInfo, error)

func StartFetchingBoardInfoUpdates(
	dvachClient *dvach.Client,
	storage *storage.Storage,
	updatesTimeout time.Duration,
	callback BoardInfoCallback,
) {
	updatesTicker := time.NewTicker(updatesTimeout)
	for _ = range updatesTicker.C {
		fetchInfoForAllBoards(dvachClient, storage, callback)
	}
}

func fetchInfoForAllBoards(
	dvachClient *dvach.Client,
	storage *storage.Storage,
	callback BoardInfoCallback,
) {
	boardNames, err := storage.AllBoardNames()
	if err != nil {
		callback(nil, err)
		return
	}
	log.Printf("%v", boardNames)

	for _, boardName := range boardNames {
		boardInfo, err := dvachClient.BoardInfo(boardName)
		callback(boardInfo, err)
	}
}
