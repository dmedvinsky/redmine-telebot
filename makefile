build: main.go
	go build

run: build
	. ./.env && ./redminebot

DEPLOY_TO = "/home/redmine/telegrambot"

deploy: build
	scp ./redminebot redmine:$(DEPLOY_TO)/redminebot.new
	ssh redmine "supervisorctl stop telegrambot && mv $(DEPLOY_TO)/redminebot.new $(DEPLOY_TO)/redminebot && supervisorctl start telegrambot"
