package bot

import (
	"encoding/json"
	"sort"
)

type ThreadInfo struct {
	Id        string `json:"num"`
	Subject   string `json:"subject"`
	Timestamp int64  `json:"timestamp"`
}

type ByTimestamp []ThreadInfo

func (threads ByTimestamp) Len() int           { return len(threads) }
func (threads ByTimestamp) Swap(i, j int)      { threads[i], threads[j] = threads[j], threads[i] }
func (threads ByTimestamp) Less(i, j int) bool { return threads[i].Timestamp < threads[j].Timestamp }

type BoardInfo struct {
	Board   string       `json:"board"`
	Threads []ThreadInfo `json:"threads"`
}

func NewBoardInfo(jsonBoardData []byte) *BoardInfo {
	boardInfo := BoardInfo{}
	json.Unmarshal(jsonBoardData, &boardInfo)

	return &boardInfo
}

func (boardInfo *BoardInfo) threadsAfter(timestamp int64) []ThreadInfo {
	threads := make([]ThreadInfo, 0)
	for _, thread := range boardInfo.Threads {
		if thread.Timestamp > timestamp {
			threads = append(threads, thread)
		}
	}
	sort.Sort(ByTimestamp(threads))

	return threads
}
