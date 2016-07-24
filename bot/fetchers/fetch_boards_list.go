package fetchers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	boardsEndpoint = "https://2ch.hk/boards.json"
)

type BoardsListCallback func(*BoardsList, error)
type BoardsListFetcher struct {
	HttpClient *http.Client
}

func StartFetchingBoardsListUpdates(boardsFetcher *BoardsListFetcher,
	updatesTimeout time.Duration,
	callback BoardsListCallback,
) {
	boardsFetcher.fetchBoardsListFromServer(callback)

	updatesTicker := time.NewTicker(updatesTimeout)
	for _ = range updatesTicker.C {
		boardsFetcher.fetchBoardsListFromServer(callback)
	}
}

func (fetcher *BoardsListFetcher) fetchBoardsListFromServer(callback BoardsListCallback) {
	response, err := fetcher.HttpClient.Get(boardsEndpoint)
	if err != nil {
		callback(nil, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		callback(nil, errors.New(fmt.Sprintf(FetchError, response.StatusCode)))
		return
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		callback(nil, err)
		return
	}

	boardList := NewBoardsList(responseBytes)
	callback(boardList, nil)
}
