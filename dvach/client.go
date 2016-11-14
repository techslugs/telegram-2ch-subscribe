package dvach

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	HttpClient *http.Client
}

const (
	FetchError      = "Error code: %v"
	FileEndpoint    = "https://2ch.hk%s"
	BoardsEndpoint  = "https://2ch.hk/boards.json"
	ThreadsEndpoint = "https://2ch.hk/%s/threads.json"
	ThreadEndpoint  = "https://2ch.hk/%s/res/%s.json"
)

func (client *Client) BoardsList() (*BoardsList, error) {
	response, err := client.HttpClient.Get(BoardsEndpoint)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf(FetchError, response.StatusCode))
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	boardList := NewBoardsList(responseBytes)
	return boardList, nil
}

func (client *Client) BoardInfo(board string) (*BoardInfo, error) {
	url := fmt.Sprintf(ThreadsEndpoint, board)
	response, err := client.HttpClient.Get(url)
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

func (client *Client) ThreadFirstPost(board string, threadID string) (*Post, error) {
	url := fmt.Sprintf(ThreadEndpoint, board, threadID)
	response, err := client.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf(FetchError, response.StatusCode))
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	boardWithThread := NewBoardWithThreadWithPosts(responseBytes)
	return boardWithThread.ThreadPost()
}
