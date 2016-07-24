package bot

import (
	"log"
	"telegram-2ch-news-bot/bot/fetchers"
)

func (bot *Bot) saveBoardsList(boardsList *fetchers.BoardsList) {
	var err error
	for _, board := range boardsList.Boards {
		err = bot.Storage.SaveBoardDescription(board.Board, board.FullName, board.Description)

		if err != nil {
			log.Printf("Error saving board list: %s\n", err.Error())
		}
	}
}
