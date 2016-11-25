package main

import (
	"log"
	"strings"
	"time"

	"github.com/tucnak/telebot"
)


func main() {
	loadConfig()
	log.Println("Using config: ", Config)
	connectDB()

	bot := mkBot()
	messages := make(chan telebot.Message, 100)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		switch {
		case strings.HasPrefix(message.Text, "/connect"):
			msg := connect(message)
			bot.SendMessage(message.Chat, msg, nil)
		case strings.HasPrefix(message.Text, "/disconnect"):
			msg := disconnect(message)
			bot.SendMessage(message.Chat, msg, nil)
		default:
			rmApi := mkRedmine(message.Sender.ID)
			issueIds := getIssueIds(message)
			for i := range issueIds {
				msg := getIssueData(rmApi, issueIds[i])
				bot.SendMessage(message.Chat, msg, nil)
			}
		}
	}
}
