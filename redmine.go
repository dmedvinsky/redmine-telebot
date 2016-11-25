package main

import (
	"log"

	"github.com/jason0x43/go-redmine"
)

func mkRedmine(userId int) redmine.Session {
	token, err := Redis.Get(senderKey(userId)).Result()
	if err != nil {
		log.Println(err)
	}
	return redmine.OpenSession(Config.RedmineUrl, token)
}
