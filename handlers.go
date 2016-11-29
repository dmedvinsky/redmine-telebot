package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/mattn/go-redmine"
	"github.com/tucnak/telebot"
)

func connect(message telebot.Message, user User) {
	var msg string
	msg = "Please use `/connect my_redmine_token`"
	parts := strings.Fields(message.Text)
	if len(parts) == 2 {
		user.RedmineToken = parts[1]
		if user.save() {
			msg = "Connected. Now you can use the bot."
		} else {
			msg = "Error"
		}
	}
	Bot.SendMessage(message.Chat, msg, nil)
}

func disconnect(message telebot.Message, user User) {
	var msg string
	if user.delete() {
		msg = "Disconnected. Your Redmine access token has been deleted."
	} else {
		msg = "Error"
	}
	Bot.SendMessage(message.Chat, msg, nil)
}

func parseMessage(message telebot.Message, user User) {
	var issueIdRe = regexp.MustCompile(`#(?P<issue>\d+)`)
	var issueLinkRe = regexp.MustCompile(Config.RedmineUrl + `/issues/(?P<issue>\d+)/?`)
	var issueIds []int
	issueIds = append(issueIds, getIds(message.Text, issueIdRe)...)
	issueIds = append(issueIds, getIds(message.Text, issueLinkRe)...)
	for i := range issueIds {
		msg := getIssueData(user, issueIds[i])
		Bot.SendMessage(message.Chat, msg, nil)
	}
}

func getIssueData(user User, id int) (msg string) {
	rmApi := redmine.NewClient(Config.RedmineUrl, user.RedmineToken)
	issue, err := rmApi.Issue(id)
	if err != nil {
		log.Println(err)
		msg = fmt.Sprintf("Issue #%d: error accessing", id)
	} else {
		url := fmt.Sprintf("%s/issues/%d", Config.RedmineUrl, issue.Id)
		msg = fmt.Sprintf("%s #%d: %s\n%s", issue.Tracker.Name, issue.Id, issue.Subject, url)
	}
	return
}
