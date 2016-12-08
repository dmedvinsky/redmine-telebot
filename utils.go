package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-redmine"
)

func getIds(message string, pattern *regexp.Regexp) (ids []int) {
	matches := pattern.FindAllStringSubmatch(message, -1)
	for i := range matches {
		id, err := strconv.Atoi(matches[i][1])
		if err != nil {
			log.Println(err)
		}
		ids = append(ids, id)
	}
	return
}

func getIdFromKeyboard(message string) (id int) {
	parts := strings.SplitN(message, ":", 2)
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Println(err)
	}
	return
}

func projectsKeyboard(items []redmine.Project) (keyboard [][]string) {
	for i := 0; i < len(items); i += 2 {
		if i > 100 {
			break
		}
		var row [2]string
		row[0] = fmt.Sprintf("%d: %s", items[i].Id, items[i].Name)
		if i+1 < len(items) {
			row[1] = fmt.Sprintf("%d: %s", items[i+1].Id, items[i+1].Name)
		}
		keyboard = append(keyboard, row[:])
	}
	return
}

func activitiesKeyboard(items []redmine.TimeEntryActivity) (keyboard [][]string) {
	for i := 0; i < len(items); i += 2 {
		if i > 20 {
			break
		}
		var row [2]string
		row[0] = fmt.Sprintf("%d: %s", items[i].Id, items[i].Name)
		if i+1 < len(items) {
			row[1] = fmt.Sprintf("%d: %s", items[i+1].Id, items[i+1].Name)
		}
		keyboard = append(keyboard, row[:])
	}
	return
}
