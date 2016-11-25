package main

import (
	"log"
	"regexp"
	"strconv"
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
