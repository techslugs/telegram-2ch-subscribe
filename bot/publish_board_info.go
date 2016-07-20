package bot

import (
	"log"
)

func (bot *Bot) publishBoardInfo(boardInfo *BoardInfo) {
	board, err := bot.storage.BoardDetails(boardInfo.Board)
	if err != nil {
		log.Printf("Error retrieving board details: %s\n", err.Error())
		return
	}

	lastTimestamp := board.Timestamp
	threads := boardInfo.threadsAfter(lastTimestamp)
	if len(threads) == 0 {
		return
	}

	for _, thread := range threads {
		log.Println(thread.Subject, thread.Timestamp)
		bot.storage.UpdateBoardTimestamp(board.Name, thread.Timestamp)
	}
	log.Println("--------------------------------")
}
