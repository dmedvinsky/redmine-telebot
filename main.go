package main

import (
	"log"
	"strings"
	"time"

	"github.com/mattn/go-redmine"
	"github.com/tucnak/telebot"
)

func main() {
	loadConfig()
	log.Println("Using config: ", Config)
	connectDB()
	mkBot()
	redmine.DefaultLimit = 100

	messages := make(chan telebot.Message, 100)
	Bot.Listen(messages, 1*time.Second)

	for message := range messages {
		user := GetUser(message.Sender.ID)
		switch {
		case strings.HasPrefix(message.Text, "/abort"):
			fallthrough
		case strings.HasPrefix(message.Text, "/cancel"):
			abort(message, user)
		case strings.HasPrefix(message.Text, "/connect"):
			connect(message, user)
		case strings.HasPrefix(message.Text, "/disconnect"):
			disconnect(message, user)
		case strings.HasPrefix(message.Text, "/track"):
			track0(message, user)
		default:
			switch user.ChatState["state_id"] {
			case "track.1":
				track1(message, user)
			case "track.2":
				track2(message, user)
			case "track.3":
				track3(message, user)
			case "track.4":
				track4(message, user)
			default:
				parseMessage(message, user)
			}
		}
	}
}
