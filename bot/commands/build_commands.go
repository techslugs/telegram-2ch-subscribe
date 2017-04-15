package commands

func BuildCommands(botName string) []Command {
	commands := []Command{
		BuildSubscribeChannel(botName),
		BuildUnsubscribeChannel(botName),
		BuildSubscribe(botName),
		BuildUnsubscribe(botName),
		BuildSetStopWordsChannel(botName),
		BuildSetStopWords(botName),
	}
	return append(commands, BuildUsage(botName, commands))
}
