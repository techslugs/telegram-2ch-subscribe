package stop_words

import (
	"regexp"
	"strings"
)

var (
	nonwordsRegexp = regexp.MustCompile(`[^\wа-я]+`)
	spacesRegexp   = regexp.MustCompile(`\s+`)
)

func Normalize(message string) string {
	message = strings.ToLower(message)
	message = nonwordsRegexp.ReplaceAllLiteralString(message, " ")
	message = spacesRegexp.ReplaceAllLiteralString(message, " ")

	return message
}

func BuildStopwordsRegexpString(stopWords []string) string {
	regexpString := ""
	for i, word := range stopWords {
		if i != 0 {
			regexpString += `|`
		}
		regexpString += `\Q` + Normalize(word) + `\E`
	}
	return regexpString
}
