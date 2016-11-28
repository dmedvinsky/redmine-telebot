package main

import (
	"fmt"

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

func senderKey(userId int) string {
	return fmt.Sprintf("redminebot:user:%d", userId)
}
