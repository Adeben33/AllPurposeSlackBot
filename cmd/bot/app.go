package bot

import (
	"slackbot/api/handlers"
	"slackbot/internal/config"
	"slackbot/services/logger"
	"strings"
)

func init() {
	config.LoadConfig(".")
	logger.LoadLogger()
}

type log *logger.Logger

func Setup() {
	logger1 := logger.GetLogger()
	logger2 := logger.NewLogger(*logger1)
	config := config.GetConfig()

	if !strings.HasPrefix(config.SlackAppToken, "xapp-") {
		panic("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := config.SlackBotToken
	if botToken == "" {
		panic("SLACK_BOT_TOKEN must be set.\n")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		panic("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}
	handlers.SlackBotHandler(logger2, config)

}
