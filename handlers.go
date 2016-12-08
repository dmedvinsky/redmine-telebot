package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
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
		if user.Save() {
			msg = "Connected. Now you can use the bot."
		} else {
			msg = "Error"
		}
	}
	Bot.SendMessage(message.Chat, msg, nil)
}

func disconnect(message telebot.Message, user User) {
	var msg string
	if user.Delete() {
		msg = "Disconnected. Your Redmine access token has been deleted."
	} else {
		msg = "Error"
	}
	Bot.SendMessage(message.Chat, msg, nil)
}

func track0(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	rmApi := redmine.NewClient(Config.RedmineUrl, user.RedmineToken)
	projects, err := rmApi.Projects()
	if err != nil {
		log.Println(err)
	}
	keyboard := projectsKeyboard(projects)
	Bot.SendMessage(message.Chat, "Please select project.",
		&telebot.SendOptions{
			ReplyTo: message,
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      true,
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
				CustomKeyboard:  keyboard,
			},
		})
	user.SetState("state_id", "track.1")
}

func track1(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	rmApi := redmine.NewClient(Config.RedmineUrl, user.RedmineToken)
	projectId := getIdFromKeyboard(message.Text)
	if projectId == 0 {
		track0(message, user)
		return
	}
	user.SetState("project_id", projectId)
	activities, err := rmApi.TimeEntryActivites()
	if err != nil {
		log.Println(err)
	}
	keyboard := activitiesKeyboard(activities)
	Bot.SendMessage(message.Chat, "Please select activity.",
		&telebot.SendOptions{
			ReplyTo: message,
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      true,
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
				CustomKeyboard:  keyboard,
			},
		})
	user.SetState("state_id", "track.2")
}

func track2(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	activityId := getIdFromKeyboard(message.Text)
	if activityId == 0 {
		track1(message, user)
		return
	}
	user.SetState("activity_id", activityId)
	Bot.SendMessage(message.Chat, "Please enter spent hours.", nil)
	user.SetState("state_id", "track.3")
}

func track3(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	hours, err := strconv.ParseFloat(message.Text, 64)
	if err != nil {
		log.Println(err)
	}
	if hours == 0 {
		track2(message, user)
		return
	}
	user.SetState("hours", hours)
	Bot.SendMessage(message.Chat, "Please enter comment.", nil)
	user.SetState("state_id", "track.4")
}

func track4(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	rmApi := redmine.NewClient(Config.RedmineUrl, user.RedmineToken)
	var te redmine.TimeEntry
	te.ProjectId = user.ChatState["project_id"].(int)
	te.ActivityId = user.ChatState["activity_id"].(int)
	te.Hours = float32(user.ChatState["hours"].(float64))
	te.Comments = message.Text
	_, err := rmApi.CreateTimeEntry(te)
	if err != nil {
		log.Println(err)
		Bot.SendMessage(message.Chat, "Fail!", nil)
	} else {
		Bot.SendMessage(message.Chat, "Successfully logged the [wasted] time.", nil)
	}
	user.ClearState()
}

func parseMessage(message telebot.Message, user User) {
	var issueIdRe = regexp.MustCompile(`#(?P<issue>\d+)`)
	var issueLinkRe = regexp.MustCompile(Config.RedmineUrl + `/issues/(?P<issue>\d+)/?`)
	var issueIds []int
	issueIds = append(issueIds, getIds(message.Text, issueIdRe)...)
	issueIds = append(issueIds, getIds(message.Text, issueLinkRe)...)
	for _, id := range issueIds {
		Bot.SendChatAction(message.Chat, telebot.Typing)
		msg := getIssueData(user, id)
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
