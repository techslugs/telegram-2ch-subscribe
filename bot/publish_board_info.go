package bot

import (
	"html"
	"log"
)

func (bot *Bot) publishBoardInfo(boardInfo *BoardInfo) {
	board, err := bot.storage.BoardDetails(boardInfo.Board)
	if err != nil {
		log.Printf("Error retrieving board details: %s\n", err.Error())
		return
	}

	lastTimestamp := board.Timestamp
	threads := boardInfo.ThreadsAfter(lastTimestamp)
	if len(threads) == 0 {
		return
	}

	var threadURL string
	for _, thread := range threads {
		threadURL = boardInfo.ThreadUrl(thread.ID)

		log.Printf(
			"[%s] %s: %s",
			boardInfo.Board,
			html.UnescapeString(thread.Subject),
			threadURL,
		)
		for _, chatID := range board.ChatIDs {
			bot.sendMessage(chatID, threadURL)
		}

		bot.storage.UpdateBoardTimestamp(board.Name, thread.Timestamp)
	}
}
