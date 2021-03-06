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
		case strings.HasPrefix(message.Text, "/abort"),
			strings.HasPrefix(message.Text, "/cancel"):
			abort(message, user)
		case strings.HasPrefix(message.Text, "/connect"):
			connect(message, user)
		case strings.HasPrefix(message.Text, "/disconnect"):
			disconnect(message, user)
		case strings.HasPrefix(message.Text, "/comment"):
			comment0(message, user)
		case strings.HasPrefix(message.Text, "/track"):
			track0(message, user)
		case strings.HasPrefix(message.Text, "/new"),
			strings.HasPrefix(message.Text, "/inprogress"),
			strings.HasPrefix(message.Text, "/review"),
			strings.HasPrefix(message.Text, "/fixed"),
			strings.HasPrefix(message.Text, "/closed"),
			strings.HasPrefix(message.Text, "/reopened"),
			strings.HasPrefix(message.Text, "/feedback"):
			changeStatus(message, user)
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
			case "track.5":
				track5(message, user)
			case "comment.1":
				comment1(message, user)
			case "comment.2":
				comment2(message, user)
			default:
				parseMessage(message, user)
			}
		}
	}
}
