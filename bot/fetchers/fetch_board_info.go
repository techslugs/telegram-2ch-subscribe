package fetchers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"telegram-2ch-news-bot/storage"
	"time"
)

const (
	FetchError      = "Error code: %v"
	threadsEndpoint = "https://2ch.hk/%s/threads.json"
)

type BoardInfoCallback func(*BoardInfo, error)
type BoardsInfoFetcher struct {
	HttpClient *http.Client
	Storage    *storage.Storage
}

func StartFetchingBoardInfoUpdates(boardsInfoFetcher *BoardsInfoFetcher,
	updatesTimeout time.Duration,
	callback BoardInfoCallback,
) {
	updatesTicker := time.NewTicker(updatesTimeout)
	for _ = range updatesTicker.C {
		boardsInfoFetcher.fetchInfoForAllBoards(callback)
	}
}

func (fetcher *BoardsInfoFetcher) fetchInfoForAllBoards(callback BoardInfoCallback) {
	boardNames, err := fetcher.Storage.AllBoardNames()
	if err != nil {
		callback(nil, err)
		return
	}

	for _, boardName := range boardNames {
		boardInfo, err := fetcher.fetchBoardInfo(boardName)
		callback(boardInfo, err)
	}
}

func (fetcher *BoardsInfoFetcher) fetchBoardInfo(board string) (*BoardInfo, error) {
	url := fmt.Sprintf(threadsEndpoint, board)
	response, err := fetcher.HttpClient.Get(url)
	if err != nil {
		return &BoardInfo{Board: board}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return &BoardInfo{Board: board},
			errors.New(fmt.Sprintf(FetchError, response.StatusCode))
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &BoardInfo{Board: board}, err
	}

	boardInfo := NewBoardInfo(responseBytes)

	return boardInfo, nil
}
