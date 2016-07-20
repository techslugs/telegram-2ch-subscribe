package bot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	threadsEndpoint = "https://2ch.hk/%s/threads.json"
)

type BoardInfoCallback func(*BoardInfo, error)

func (bot *Bot) fetchBoardInfoUpdates(updatesTimeout time.Duration, callback BoardInfoCallback) {
	updatesTicker := time.NewTicker(updatesTimeout)
	for _ = range updatesTicker.C {
		bot.fetchInfoForAllBoards(callback)
	}
}

func (bot *Bot) fetchInfoForAllBoards(callback BoardInfoCallback) {
	boardInfo, err := bot.fetchBoardInfo("b")
	callback(boardInfo, err)
}

func (bot *Bot) fetchBoardInfo(board string) (*BoardInfo, error) {
	url := fmt.Sprintf(threadsEndpoint, board)
	httpClient := http.Client{}
	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	boardInfo := NewBoardInfo(responseBytes)

	return boardInfo, nil
}
