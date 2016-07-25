package fetchers

import (
	"github.com/tmwh/telegram-2ch-subscribe/dvach"
	"time"
)

const (
	boardsEndpoint = "https://2ch.hk/boards.json"
)

type BoardsListCallback func(*dvach.BoardsList, error)

func StartFetchingBoardsListUpdates(
	dvachClient *dvach.Client,
	updatesTimeout time.Duration,
	callback BoardsListCallback,
) {
	boardsList, err := dvachClient.BoardsList()
	callback(boardsList, err)

	updatesTicker := time.NewTicker(updatesTimeout)
	for _ = range updatesTicker.C {
		boardsList, err := dvachClient.BoardsList()
		callback(boardsList, err)
	}
}
