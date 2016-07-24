package commands

import (
	"regexp"
)

var Usage = &BaseCommand{
	regexp: regexp.MustCompile(`\s*/2ch_usage`),
	successMessage: Subscribe.UsageMessage() +
		"\n" +
		Unsubscribe.UsageMessage() +
		"\n" +
		SubscribeChannel.UsageMessage() +
		"\n" +
		UnsubscribeChannel.UsageMessage(),
}
