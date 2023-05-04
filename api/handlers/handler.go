package handlers

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"slackbot/internal/config"
	"slackbot/middleware"
	"slackbot/services/logger"
)

func SlackBotHandler(logger *logger.Logger, config *config.Config) {
	api := slack.New(config.SlackBotToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(config.SlackAppToken))

	client := socketmode.New(api,
		socketmode.OptionDebug(true))

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	socketmodeHandler.Handle(socketmode.EventTypeConnecting, middleware.MiddlewareConnecting)
	socketmodeHandler.Handle(socketmode.EventTypeConnectionError, middleware.MiddlewareConnectionError)
	socketmodeHandler.Handle(socketmode.EventTypeConnected, middleware.MiddlewareConnected)

	// Handle all EventsAPI
	socketmodeHandler.Handle(socketmode.EventTypeEventsAPI, middleware.MiddlewareEventsAPI)

	// Handle a specific event from EventsAPI
	//socketmodeHandler.HandleEvents(slackevents.AppMention, middlewareAppMentionEvent)

	//slash

	socketmodeHandler.HandleSlashCommand("/age", middleware.MiddlewareSlashAge)

	socketmodeHandler.RunEventLoop()
	logger.Info("Starting")
}
