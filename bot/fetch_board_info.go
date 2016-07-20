package bot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	FetchBoardInfoError = "Error code: %v"
	threadsEndpoint     = "https://2ch.hk/%s/threads.json"
)

type BoardInfoCallback func(*BoardInfo, error)

func (bot *Bot) fetchBoardInfoUpdates(updatesTimeout time.Duration, callback BoardInfoCallback) {
	updatesTicker := time.NewTicker(updatesTimeout)
	for _ = range updatesTicker.C {
		bot.fetchInfoForAllBoards(callback)
	}
}

func (bot *Bot) fetchInfoForAllBoards(callback BoardInfoCallback) {
	boardNames, err := bot.storage.AllBoardNames()
	if err != nil {
		callback(nil, err)
		return
	}

	for _, boardName := range boardNames {
		boardInfo, err := bot.fetchBoardInfo(boardName)
		callback(boardInfo, err)
	}
}

func (bot *Bot) fetchBoardInfo(board string) (*BoardInfo, error) {
	url := fmt.Sprintf(threadsEndpoint, board)
	httpClient := http.Client{}
	response, err := httpClient.Get(url)
	if err != nil {
		return &BoardInfo{Board: board}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return &BoardInfo{Board: board},
			errors.New(fmt.Sprintf(FetchBoardInfoError, response.StatusCode))
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &BoardInfo{Board: board}, err
	}

	boardInfo := NewBoardInfo(responseBytes)

	return boardInfo, nil
}
