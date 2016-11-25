package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jason0x43/go-redmine"
	"github.com/tucnak/telebot"
)

var issueIdRe = regexp.MustCompile(`#(?P<issue>\d+)`)
var issueLinkRe = regexp.MustCompile(Config.RedmineUrl + `/issues/(?P<issue>\d+)/?`)

func connect(message telebot.Message) string {
	parts := strings.Fields(message.Text)
	if len(parts) == 2 {
		err := Redis.Set(senderKey(message.Sender.ID), parts[1], 0).Err()
		if err != nil {
			log.Println(err)
			return "Error"
		}
		return "Connected. Now you can use the bot."
	}
	return "Please use `/connect my_redmine_token`"
}

func disconnect(message telebot.Message) string {
	err := Redis.Del(senderKey(message.Sender.ID)).Err()
	if err != nil {
		log.Println(err)
		return "Error"
	}
	return "Disconnected. Your Redmine access token has been deleted."
}

func getIssueIds(message telebot.Message) (ids []int) {
	ids = append(ids, getIds(message.Text, issueIdRe)...)
	ids = append(ids, getIds(message.Text, issueLinkRe)...)
	return
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
