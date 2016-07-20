package bot

import (
	"fmt"
	"time"
)

func (bot *Bot) publishBoardInfo(boardInfo *BoardInfo) {
	lastTimestamp := lastPublishedTimestamp(boardInfo)
	threads := boardInfo.threadsAfter(lastTimestamp)
	if len(threads) == 0 {
		return
	}

	for _, thread := range threads {
		fmt.Println(thread.Subject)
	}
	fmt.Println("--------------------------------")
}

func lastPublishedTimestamp(boardInfo *BoardInfo) int64 {
	now := time.Now().Add(-5 * time.Minute)
	return now.Unix()
}
