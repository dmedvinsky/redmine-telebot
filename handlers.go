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

var issueIdRe = regexp.MustCompile(`#(?P<issue>\d+)`)
var issueLinkRe = regexp.MustCompile(Config.RedmineUrl + `/issues/(?P<issue>\d+)/?`)

func abort(message telebot.Message, user User) {
	user.ClearState()
	Bot.SendMessage(message.Chat, "Cancelled current operation.",
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				HideCustomKeyboard: true,
			},
		})
}

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
	filter := redmine.IssueFilter{
		AssignedToId: "me",
		ProjectId:    strconv.Itoa(projectId),
		StatusId:     "open",
	}
	issues, err := rmApi.IssuesByFilter(&filter)
	if err != nil {
		log.Println(err)
	}
	keyboard := issuesKeyboard(issues)
	Bot.SendMessage(message.Chat, "Please select issue.",
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
	rmApi := redmine.NewClient(Config.RedmineUrl, user.RedmineToken)
	if message.Text != "No Issue" {
		issueId := getIdFromKeyboard(message.Text)
		user.SetState("issue_id", issueId)
	}
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
	user.SetState("state_id", "track.3")
}

func track3(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	activityId := getIdFromKeyboard(message.Text)
	if activityId == 0 {
		track0(message, user)
		return
	}
	user.SetState("activity_id", activityId)
	Bot.SendMessage(message.Chat, "Please enter spent hours.",
		&telebot.SendOptions{
			ReplyTo: message,
			ReplyMarkup: telebot.ReplyMarkup{
				HideCustomKeyboard: true,
			},
		})
	user.SetState("state_id", "track.4")
}

func track4(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	hours, _ := strconv.ParseFloat(message.Text, 64)
	if hours == 0 {
		track0(message, user)
		return
	}
	user.SetState("hours", hours)
	Bot.SendMessage(message.Chat, "Please enter comment.",
		&telebot.SendOptions{
			ReplyTo: message,
			ReplyMarkup: telebot.ReplyMarkup{
				HideCustomKeyboard: true,
			},
		})
	user.SetState("state_id", "track.5")
}

func track5(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	rmApi := redmine.NewClient(Config.RedmineUrl, user.RedmineToken)
	var te redmine.TimeEntry
	te.ProjectId = user.ChatState["project_id"].(int)
	te.ActivityId = user.ChatState["activity_id"].(int)
	te.Hours = float32(user.ChatState["hours"].(float64))
	te.Comments = message.Text
	issueId, ok := user.ChatState["issue_id"]
	if ok {
		te.IssueId = issueId.(int)
	}
	_, err := rmApi.CreateTimeEntry(te)
	var msg string
	if err != nil {
		log.Println(err)
		msg = "Fail!"
	} else {
		msg = "Successfully logged the [wasted] time."
	}
	Bot.SendMessage(message.Chat, msg,
		&telebot.SendOptions{
			ReplyTo: message,
			ReplyMarkup: telebot.ReplyMarkup{
				HideCustomKeyboard: true,
			},
		})
	user.ClearState()
}

func changeStatus(message telebot.Message, user User) {
	Bot.SendChatAction(message.Chat, telebot.Typing)
	rmApi := redmine.NewClient(Config.RedmineUrl, user.RedmineToken)
	statusMap := map[string]string{
		"new":        "New",
		"inprogress": "In Progress",
		"review":     "CodeReview",
		"fixed":      "Fixed",
		"closed":     "Closed",
		"reopened":   "ReOpened",
		"feedback":   "Feedback",
	}
	requestedStatus := statusMap[strings.Fields(message.Text)[0][1:]]
	issueIds := getIds(message.Text, issueIdRe)
	if issueIds == nil {
		Bot.SendMessage(message.Chat, "Please provide at least one Issue ID", nil)
		return
	}
	statuses, err := rmApi.IssueStatuses()
	if err != nil {
		log.Println(err)
		Bot.SendMessage(message.Chat, "Error accessing Redmine", nil)
		return
	}
	var status redmine.IssueStatus
	for _, s := range statuses {
		if strings.EqualFold(s.Name, requestedStatus) {
			status = s
			break
		}
	}
	if status.Id == 0 {
		Bot.SendMessage(message.Chat, "Status not found", nil)
		return
	}
	for _, id := range issueIds {
		issue, err := rmApi.Issue(id)
		if err != nil {
			log.Println(err)
			msg := fmt.Sprintf("Issue #%d: error accessing", id)
			Bot.SendMessage(message.Chat, msg, nil)
			continue
		}
		issue.ProjectId = issue.Project.Id
		issue.TrackerId = issue.Tracker.Id
		issue.StatusId = status.Id
		err = rmApi.UpdateIssue(*issue)
		if err != nil {
			log.Println(err)
			msg := fmt.Sprintf("Issue #%d: error updating", id)
			Bot.SendMessage(message.Chat, msg, nil)
			continue
		}
		msg := fmt.Sprintf("Updated #%d: %s\n%s", issue.Id, issue.Subject, issueUrl(issue))
		Bot.SendMessage(message.Chat, msg, nil)
	}
}

func parseMessage(message telebot.Message, user User) {
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
		msg = fmt.Sprintf("%s #%d: %s (%s)\n%s", issue.Tracker.Name, issue.Id,
			issue.Subject, issue.Status.Name, issueUrl(issue))
	}
	return
}
