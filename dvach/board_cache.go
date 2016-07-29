package dvach

import (
	"fmt"
	"strings"
	"sync"
)

type BoardCache struct {
	DvachClient *Client
	BoardName   string
	threads     map[string]string
	sync.RWMutex
}

func NewBoardCache(dvachClient *Client, boardName string) *BoardCache {
	return &BoardCache{
		DvachClient: dvachClient,
		BoardName:   boardName,
		threads:     make(map[string]string),
	}
}

func (cache *BoardCache) GetFormattedThreadMessage(threadID string) (string, error) {
	cache.RLock()
	cachedMessage, ok := cache.threads[threadID]
	cache.RUnlock()
	if ok {
		return cachedMessage, nil
	}

	threadMessage, err := getFormattedThreadMessage(cache.DvachClient, cache.BoardName, threadID)
	if err != nil {
		return threadMessage, err
	}

	cache.Lock()
	cache.threads[threadID] = threadMessage
	cache.Unlock()
	return threadMessage, nil
}

func getFormattedThreadMessage(
	dvachClient *Client,
	board, threadID string,
) (string, error) {
	post, err := dvachClient.ThreadFirstPost(board, threadID)
	if err != nil {
		return "", err
	}

	message := ""
	if post.Subject != "" {
		message = message + fmt.Sprintf("*%s*\n", post.SanitizedSubject())
	}
	if fileURL := post.FileUrl(board); fileURL != "" {
		message = message + fmt.Sprintf("%s\n\n", fileURL)
	}
	if comment := post.SanitizedComment(); comment != "" {
		message = message + fmt.Sprintf("%s\n", comment)
	}
	message = fmt.Sprintf("%.4000s", message)
	message = addMissingFormatting(message, "*")
	message = addMissingFormatting(message, "_")
	message = message + post.ThreadUrl(board)

	return message, nil
}

func addMissingFormatting(message, formatChar string) string {
	allCount := strings.Count(message, formatChar)
	escapedCount := strings.Count(message, "\\"+formatChar)
	if (allCount-escapedCount)%2 == 0 {
		return message
	}

	return message + formatChar
}
