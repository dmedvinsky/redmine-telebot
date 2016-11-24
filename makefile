build: main.go
	go build

run:
	. ./.env && go run main.go

DEPLOY_TO = "/home/redmine/telegrambot"

deploy: build
	scp ./redminebot redmine:$(DEPLOY_TO)/redminebot.new
	ssh redmine "supervisorctl stop telegrambot && mv $(DEPLOY_TO)/redminebot.new $(DEPLOY_TO)/redminebot && supervisorctl start telegrambot"
