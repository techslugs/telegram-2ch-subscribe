package fetchers

import (
	"encoding/json"
)

type BoardsList struct {
	Boards []BoardDescription `json:"boards"`
}

type BoardDescription struct {
	Board       string `json:"id"`
	FullName    string `json:"name"`
	Description string `json:"info"`
}

func NewBoardsList(jsonBoardsData []byte) *BoardsList {
	boardsList := BoardsList{}
	json.Unmarshal(jsonBoardsData, &boardsList)

	return &boardsList
}
