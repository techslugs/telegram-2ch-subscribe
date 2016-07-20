package bot

import (
	"encoding/json"
)

type ThreadUpdate struct {
	id        string `json:"num"`
	subject   string `json:"subject"`
	timestamp int64  `json:"timestamp"`
	board     string `json:"board"`
}

func NewThreadUpdate(jsonThreadData []byte) *ThreadUpdate {
	threadUpdate := ThreadUpdate{}
	json.Unmarshal(jsonThreadData, &threadUpdate)

	return &threadUpdate
}
