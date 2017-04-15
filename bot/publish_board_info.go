package bot

import (
	"fmt"
	"github.com/techslugs/telegram-2ch-subscribe/dvach"
	"github.com/techslugs/telegram-2ch-subscribe/storage"
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

	bot.clearStaleThreadIDs(boardInfo)
}

func (bot *Bot) clearStaleThreadIDs(boardInfo *dvach.BoardInfo) {
	threadIDs := boardInfo.ThreadIDs()
	if len(threadIDs) < 1 {
		return
	}

	err := bot.Storage.ClearStaleThreadIDs(boardInfo.Board, threadIDs)
	if err != nil {
		log.Printf("Error clearing stale threads IDs: %s\n", err)
	}
}

func (bot *Bot) publishThreadsToSubscription(
	boardCache *dvach.BoardCache,
	boardInfo *dvach.BoardInfo,
	boardSubscription *storage.BoardSubscription,
) {
	threads := getNotSentThreads(boardInfo, boardSubscription)
	if len(threads) == 0 {
		return
	}

	var threadURL, threadSubject string
	for _, thread := range threads {
		threadURL = boardInfo.ThreadUrl(thread.ID)
		threadSubject = html.UnescapeString(thread.Subject)

		logTagged(boardSubscription, fmt.Sprintf("%s: %s", threadSubject, threadURL))

		threadMessage, err := boardCache.GetFormattedThreadMessage(thread.ID)

		if err != nil {
			log.Printf("Error: Could not get formatted message. %s", err)
			threadMessage = threadURL
		}

		if !boardSubscription.HasStopWords(threadMessage) {
			bot.TelegramClient.SendMarkdownMessage(boardSubscription.ChatID, threadMessage)
		} else {
			logTagged(boardSubscription, fmt.Sprintf("Skipping thread %s due to stop words.", threadURL))
		}

		bot.Storage.LogSentThread(
			boardSubscription.BoardName,
			boardSubscription.ChatID,
			thread.ID,
		)
	}
}

func logTagged(boardSubscription *storage.BoardSubscription, message string) {
	log.Printf(
		"[%v %s] %s",
		boardSubscription.ChatID,
		boardSubscription.BoardName,
		message,
	)
}

func getNotSentThreads(
	boardInfo *dvach.BoardInfo,
	boardSubscription *storage.BoardSubscription,
) []dvach.ThreadInfo {
	sentThreadIDsMap := buildThreadIDsMap(boardSubscription.SentThreadIDs)

	threads := boardInfo.FilteredThreads(func(thread *dvach.ThreadInfo) bool {
		_, ok := sentThreadIDsMap[thread.ID]
		return !ok &&
			thread.Score >= boardSubscription.MinScore &&
			thread.Timestamp >= boardSubscription.Timestamp
	})
	return threads
}

func buildThreadIDsMap(sentThreadIDs []string) map[string]struct{} {
	threadIDsMap := make(map[string]struct{})
	for _, threadID := range sentThreadIDs {
		threadIDsMap[threadID] = struct{}{}
	}
	return threadIDsMap
}
