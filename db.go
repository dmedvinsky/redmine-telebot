package main

import (
	"fmt"
	"log"

	"gopkg.in/redis.v5"
)

var Redis *redis.Client

func connectDB() {
	Redis = redis.NewClient(&redis.Options{
		Network: Config.RedisNetwork,
		Addr:    Config.RedisAddr,
		DB:      Config.RedisDB,
	})
}

type User struct {
	Id           int
	RedmineToken string
}

func (u *User) dbKey() string {
	return fmt.Sprintf("redminebot:user:%d", u.Id)
}

func (u *User) load() {
	token := Redis.Get(u.dbKey()).Val()
	u.RedmineToken = token
}

func (u *User) save() bool {
	err := Redis.Set(u.dbKey(), u.RedmineToken, 0).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *User) delete() bool {
	err := Redis.Del(u.dbKey()).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func getUser(userId int) User {
	user := User{Id: userId}
	user.load()
	return user
}
