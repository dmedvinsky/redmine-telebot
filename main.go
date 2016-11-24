package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jason0x43/go-redmine"
	"github.com/tucnak/telebot"
)

var redmineUrl = os.Getenv("REDMINE_URL")
var redmineToken = os.Getenv("REDMINE_TOKEN")
var botToken = os.Getenv("BOT_TOKEN")

func main() {
	bot := mkBot()
	rmApi := mkRedmine()

	messages := make(chan telebot.Message, 100)
	bot.Listen(messages, 1*time.Second)

	issueIdRe := regexp.MustCompile(`#(?P<issue>\d+)`)
	issueLinkRe := regexp.MustCompile(redmineUrl + `/issues/(?P<issue>\d+)/?`)

	for message := range messages {
		var issueIds []int
		issueIds = append(issueIds, getIds(message.Text, issueIdRe)...)
		issueIds = append(issueIds, getIds(message.Text, issueLinkRe)...)

		for i := range issueIds {
			issue, err := rmApi.GetIssue(issueIds[i])
			var msg string
			if err != nil {
				log.Println(err)
				msg = fmt.Sprintf("Issue #%d: error accessing", issueIds[i])
			} else {
				url := fmt.Sprintf("%s/issues/%d", redmineUrl, issue.Id)
				msg = fmt.Sprintf("%s #%d: %s\n%s", issue.Tracker.Name, issue.Id, issue.Subject, url)
			}
			bot.SendMessage(message.Chat, msg, nil)
		}
	}
}

func mkBot() *telebot.Bot {
	bot, err := telebot.NewBot(botToken)
	if err != nil {
		log.Fatalln(err)
	}
	return bot
}

func mkRedmine() redmine.Session {
	return redmine.OpenSession(redmineUrl, redmineToken)
}

func getIds(message string, pattern *regexp.Regexp) []int {
	var issueIds []int
	matches := pattern.FindAllStringSubmatch(message, -1)
	if matches == nil {
		return nil
	}
	for i := range matches {
		issueId, err := strconv.Atoi(matches[i][1])
		if err != nil {
			log.Fatalln(err)
		}
		issueIds = append(issueIds, issueId)
	}
	return issueIds
}
