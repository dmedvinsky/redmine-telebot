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
		user := getUser(message.Sender.ID)
		switch {
		case strings.HasPrefix(message.Text, "/connect"):
			connect(message, user)
		case strings.HasPrefix(message.Text, "/disconnect"):
			disconnect(message, user)
		default:
			parseMessage(message, user)
		}
	}
}
