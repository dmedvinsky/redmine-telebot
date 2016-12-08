package main

import (
	"fmt"
	"log"

	"gopkg.in/redis.v5"
)

var Redis *redis.Client

var ChatStates map[int]map[string]interface{}

func connectDB() {
	Redis = redis.NewClient(&redis.Options{
		Network: Config.RedisNetwork,
		Addr:    Config.RedisAddr,
		DB:      Config.RedisDB,
	})
	ChatStates = make(map[int]map[string]interface{})
}

type User struct {
	Id           int
	RedmineToken string
	ChatState    map[string]interface{}
}

func (u *User) dbKey() string {
	return fmt.Sprintf("redminebot:user:%d", u.Id)
}

func (u *User) load() {
	token := Redis.Get(u.dbKey()).Val()
	u.RedmineToken = token
	state, ok := ChatStates[u.Id]
	if ok {
		u.ChatState = state
	} else {
		u.ChatState = make(map[string]interface{})
	}
}

func (u *User) Save() bool {
	err := Redis.Set(u.dbKey(), u.RedmineToken, 0).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *User) Delete() bool {
	err := Redis.Del(u.dbKey()).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *User) SetState(key string, val interface{}) {
	u.ChatState[key] = val
	ChatStates[u.Id] = u.ChatState
}

func (u *User) ClearState() {
	u.ChatState = make(map[string]interface{})
	ChatStates[u.Id] = u.ChatState
}

func GetUser(userId int) User {
	user := User{Id: userId}
	user.load()
	return user
}
