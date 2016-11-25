package main

import (
	"github.com/kelseyhightower/envconfig"
)

type ConfigSpec struct {
	RedisNetwork string `envconfig:"redis_network" default:"tcp"`
	RedisAddr    string `envconfig:"redis_addr" default:"localhost:6379"`
	RedisDB      int    `envconfig:"redis_db" default:"0"`
	RedmineUrl   string `envconfig:"redmine_url"`
	BotToken     string `envconfig:"bot_token"`
}

var Config ConfigSpec

func loadConfig() {
	err := envconfig.Process("redminebot", &Config)
	if err != nil {
		panic(err)
	}
}
