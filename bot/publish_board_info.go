package bot

import (
	"fmt"
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

	var threadURL, threadSubject string
	for _, thread := range threads {
		threadURL = boardInfo.ThreadUrl(thread.ID)
		threadSubject = html.UnescapeString(thread.Subject)

		log.Printf(
			"[%s] %s: %s",
			boardInfo.Board,
			threadSubject,
			threadURL,
		)

		threadMessage, err := getFormatedThreadMessage(
			bot.DvachClient,
			boardInfo.Board,
			thread.ID,
		)
		if err != nil {
			log.Printf("Error: Could not get formatted message. %s", err)
			threadMessage = threadURL
		}

		for _, chatID := range board.ChatIDs {
			bot.TelegramClient.SendMarkdownMessage(chatID, threadMessage)
		}

		bot.Storage.UpdateBoardTimestamp(board.Name, thread.Timestamp)
	}
}

func getFormatedThreadMessage(
	dvachClient *dvach.Client,
	board, threadID string,
) (string, error) {
	post, err := dvachClient.ThreadFirstPost(board, threadID)
	if err != nil {
		return "", err
	}

	message := ""
	if post.Subject != "" {
		message = message + fmt.Sprintf("*%s*\n", html.UnescapeString(post.Subject))
	}
	if fileURL := post.FileUrl(board); fileURL != "" {
		message = message + fmt.Sprintf("%s\n\n", fileURL)
	}
	if comment := post.SanitizedComment(); comment != "" {
		message = message + fmt.Sprintf("%s\n", comment)
	}
	return message + post.ThreadUrl(board), nil
}
