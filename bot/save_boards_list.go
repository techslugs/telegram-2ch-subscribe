package bot

import (
	"github.com/tmwh/telegram-2ch-subscribe/dvach"
	"log"
)

func (bot *Bot) saveBoardsList(boardsList *dvach.BoardsList) {
	var err error
	for _, board := range boardsList.Boards {
		err = bot.Storage.SaveBoardDescription(board.Board, board.FullName, board.Description)

		if err != nil {
			log.Printf("Error saving board list: %s\n", err)
		}
	}
}
