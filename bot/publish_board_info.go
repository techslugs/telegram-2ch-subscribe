package bot

import (
	"github.com/tmwh/telegram-2ch-subscribe/dvach"
	"github.com/tmwh/telegram-2ch-subscribe/storage"
	"html"
	"log"
	"sync"
)

func (bot *Bot) publishBoardInfo(boardInfo *dvach.BoardInfo) {
	boardSubscriptions, err := bot.Storage.AllBoardSubscriptions(boardInfo.Board)
	if err != nil {
		log.Printf("Error retrieving board subscriptions: %s\n", err)
		return
	}
	boardCache := dvach.NewBoardCache(bot.DvachClient, boardInfo.Board)

	wg := sync.WaitGroup{}
	wg.Add(len(boardSubscriptions))
	for _, boardSubscription := range boardSubscriptions {
		go func(boardSubscription storage.BoardSubscription) {
			bot.publishThreadsToSubscription(
				boardCache,
				boardInfo,
				&boardSubscription,
			)
			wg.Done()
		}(boardSubscription)
	}
	wg.Wait()
}

func (bot *Bot) publishThreadsToSubscription(
	boardCache *dvach.BoardCache,
	boardInfo *dvach.BoardInfo,
	boardSubscription *storage.BoardSubscription,
) {
	threads := boardInfo.NotSentThreadsWithScoreGreaterThan(
		boardSubscription.SentThreadIDs,
		boardSubscription.Timestamp,
		boardSubscription.MinScore,
	)
	if len(threads) == 0 {
		return
	}

	var threadURL, threadSubject string
	for _, thread := range threads {
		threadURL = boardInfo.ThreadUrl(thread.ID)
		threadSubject = html.UnescapeString(thread.Subject)

		log.Printf(
			"%v [%s] %s: %s",
			boardSubscription.ChatID,
			boardInfo.Board,
			threadSubject,
			threadURL,
		)

		threadMessage, err := boardCache.GetFormattedThreadMessage(thread.ID)

		if err != nil {
			log.Printf("Error: Could not get formatted message. %s", err)
			threadMessage = threadURL
		}

		bot.TelegramClient.SendMarkdownMessage(boardSubscription.ChatID, threadMessage)
		bot.Storage.LogSentThread(
			boardSubscription.BoardName,
			boardSubscription.ChatID,
			thread.ID,
		)
	}
}
