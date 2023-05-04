package middleware

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"strings"
	"time"
)

func MiddlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connecting to Slack with Socket Mode...")
}

func MiddlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connected to Slack with Socket Mode.")
}

func MiddlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connection failed. Retrying later...")
}

func MiddlewareEventsAPI(evt *socketmode.Event, client *socketmode.Client) {
	log.Println("middlewareEventsAPI")
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	fmt.Printf("Event received: %+v\n", eventsAPIEvent)

	client.Ack(*evt.Request)

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			fmt.Printf("We have been mentionned in %v", ev.Channel)
			_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				fmt.Printf("failed posting message: %v", err)
			}
		case *slackevents.MemberJoinedChannelEvent:
			fmt.Printf("user %q joined to channel %q", ev.User, ev.Channel)
		}
	default:
		client.Debugf("unsupported Events API event received")
	}
}

func MiddlewareSlashAge(evt *socketmode.Event, client *socketmode.Client) {
	cmd, ok := evt.Data.(slack.SlashCommand)

	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	client.Debugf("Slash command received: %+v", cmd)
	// Get the text of the slash command and parse it as a date
	dateStr := strings.TrimSpace(cmd.Text)
	date, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		payload := map[string]interface{}{
			"response_type": "ephemeral",
			"text":          fmt.Sprintf("Invalid date format: %s. Please use YYYY-MM-DD.", dateStr),
		}
		client.Ack(*evt.Request, payload)
		return
	}

	// Calculate the age based on the input date
	now := time.Now()
	age := now.Year() - date.Year()
	if now.Month() < date.Month() || (now.Month() == date.Month() && now.Day() < date.Day()) {
		age--
	}
	// Send the age as a response back to Slack
	payload := map[string]interface{}{
		"response_type": "in_channel",
		"text":          fmt.Sprintf("You are %d years old.", age),
	}
	client.Ack(*evt.Request, payload)
}
