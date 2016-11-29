package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jason0x43/go-redmine"
	"github.com/tucnak/telebot"
)

func connect(message telebot.Message) {
	var msg string
	parts := strings.Fields(message.Text)
	if len(parts) == 2 {
		err := Redis.Set(senderKey(message.Sender.ID), parts[1], 0).Err()
		if err != nil {
			log.Println(err)
			msg = "Error"
		}
		msg = "Connected. Now you can use the bot."
	}
	msg = "Please use `/connect my_redmine_token`"
	Bot.SendMessage(message.Chat, msg, nil)
}

func disconnect(message telebot.Message) {
	var msg string
	err := Redis.Del(senderKey(message.Sender.ID)).Err()
	if err != nil {
		log.Println(err)
		msg = "Error"
	}
	msg = "Disconnected. Your Redmine access token has been deleted."
	Bot.SendMessage(message.Chat, msg, nil)
}

func parseMessage(message telebot.Message, rmApi redmine.Session) {
	var issueIdRe = regexp.MustCompile(`#(?P<issue>\d+)`)
	var issueLinkRe = regexp.MustCompile(Config.RedmineUrl + `/issues/(?P<issue>\d+)/?`)
	var issueIds []int
	issueIds = append(issueIds, getIds(message.Text, issueIdRe)...)
	issueIds = append(issueIds, getIds(message.Text, issueLinkRe)...)
	for i := range issueIds {
		msg := getIssueData(rmApi, issueIds[i])
		Bot.SendMessage(message.Chat, msg, nil)
	}
}

func getIssueData(rmApi redmine.Session, id int) (msg string) {
	issue, err := rmApi.GetIssue(id)
	if err != nil {
		log.Println(err)
		msg = fmt.Sprintf("Issue #%d: error accessing", id)
	} else {
		url := fmt.Sprintf("%s/issues/%d", Config.RedmineUrl, issue.Id)
		msg = fmt.Sprintf("%s #%d: %s\n%s", issue.Tracker.Name, issue.Id, issue.Subject, url)
	}
	return
}
