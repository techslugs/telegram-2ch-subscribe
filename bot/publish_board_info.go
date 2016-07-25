package bot

import (
	"github.com/tmwh/telegram-2ch-subscribe/dvach"
	"html"
	"log"
)

func (bot *Bot) publishBoardInfo(boardInfo *dvach.BoardInfo) {
	board, err := bot.Storage.BoardDetails(boardInfo.Board)
	if err != nil {
		log.Printf("Error retrieving board details: %s\n", err)
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
			bot.TelegramClient.SendMessage(chatID, threadURL)
		}

		bot.Storage.UpdateBoardTimestamp(board.Name, thread.Timestamp)
	}
}
