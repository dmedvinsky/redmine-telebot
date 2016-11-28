package main

import (
	"github.com/jason0x43/go-redmine"
)

func mkRedmine(userId int) redmine.Session {
	token := Redis.Get(senderKey(userId)).Val()
	return redmine.OpenSession(Config.RedmineUrl, token)
}
