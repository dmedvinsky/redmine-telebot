package main

import (
	"log"

	"github.com/tucnak/telebot"
)

var Bot *telebot.Bot

func mkBot() {
	bot, err := telebot.NewBot(Config.BotToken)
	if err != nil {
		log.Fatalln(err)
	}
	Bot = bot
}
