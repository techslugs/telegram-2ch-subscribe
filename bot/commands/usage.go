package commands

import (
	"regexp"
	"fmt"
)

func BuildUsage(botName string, commands []Command) Command {
	regexp_template := `\s*/usage(?:@%s)?`
	regexp_source := fmt.Sprintf(regexp_template, botName)
	return &BaseCommand{
		regexp: regexp.MustCompile(regexp_source),
		successMessage: buildUsageMessage(commands),
	}
}

func buildUsageMessage(commands []Command) string {
	message := ""
	for _, command := range commands {
		message += command.UsageMessage() + "\n\n"
	}
	return message
}
