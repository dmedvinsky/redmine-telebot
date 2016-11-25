package main

import (
	"log"

	"github.com/tucnak/telebot"
)

func mkBot() (bot *telebot.Bot) {
	bot, err := telebot.NewBot(Config.BotToken)
	if err != nil {
		log.Fatalln(err)
	}
	return
}
