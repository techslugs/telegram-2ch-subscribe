package dvach

import (
	"encoding/json"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"html"
	"strconv"
	"strings"
)

type BoardWithThreadWithPosts struct {
	Board   string            `json:"board"`
	Threads []ThreadWithPosts `json:"threads"`
}

type ThreadWithPosts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID      int    `json:"num"`
	Comment string `json:"comment"`
	Subject string `json:"subject"`
	Files   []File `json:"files"`
}

type File struct {
	Path      string `json:"path"`
	Thumbnail string `json:"thumbnail"`
}

func NewBoardWithThreadWithPosts(jsonBoardData []byte) *BoardWithThreadWithPosts {
	boardWithThread := BoardWithThreadWithPosts{}
	json.Unmarshal(jsonBoardData, &boardWithThread)

	return &boardWithThread
}

func (boardWithThread *BoardWithThreadWithPosts) ThreadPost() *Post {
	return &boardWithThread.Threads[0].Posts[0]
}

func (post *Post) ThreadUrl(board string) string {
	return fmt.Sprintf(ThreadUrl, board, strconv.Itoa(post.ID))
}

func (post *Post) FileUrl(board string) string {
	if len(post.Files) == 0 {
		return ""
	}
	return fmt.Sprintf(FileEndpoint, board, post.Files[0].Path)
}

func (post *Post) SanitizedComment() string {
	comment := html.UnescapeString(post.Comment)
	comment = strings.Replace(comment, "<br>", "\n", -1)
	comment = strings.Replace(comment, "<strong>", "*", -1)
	comment = strings.Replace(comment, "</strong>", "*", -1)
	comment = strings.Replace(comment, "<i>", "_", -1)
	comment = strings.Replace(comment, "</i>", "_", -1)
	policy := bluemonday.StrictPolicy()
	return policy.Sanitize(comment)
}
