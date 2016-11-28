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
	mkBot()

	messages := make(chan telebot.Message, 100)
	Bot.Listen(messages, 1*time.Second)

	for message := range messages {
		rmApi := mkRedmine(message.Sender.ID)
		switch {
		case strings.HasPrefix(message.Text, "/connect"):
			connect(message)
		case strings.HasPrefix(message.Text, "/disconnect"):
			disconnect(message)
		default:
			parseMessage(message, rmApi)
		}
	}
}
