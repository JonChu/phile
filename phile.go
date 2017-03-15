package main

import (
	"os"
	"strings"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/gonduit/core"
	"github.com/nlopes/slack"
    "github.com/joho/godotenv"

	log "github.com/Sirupsen/logrus"
)

var botId string

func run(slackClient *slack.Client, phabClient *gonduit.Conn) int {
	log.Info("Setting up RTM connection...")

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	//log.Info("Connection successfully established.")
	log.Info("Accepting incoming events...")

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			//log.Info("Incoming event received...\n")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID
				log.Info("Connection established...")
				log.Infof("BotId: %v", botId)

			case *slack.HelloEvent:
				log.Info("Server ack received...")

			case *slack.MessageEvent:
				log.Infof("Message: %v", ev)

				if isForBot(ev) {
					/* parseMessage(ev)
					   if isCommand(parts) {
					       handleCommand()
					   } else if isResponse(parts) {

					   } else {
					       rtm.SendMessage(rtm.NewOutgoingMessage("", ev.Channel))
					   } */

				}

			case *slack.RTMError:
				log.Error("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Error("Invalid credentials")
				return 1
			}
		}
	}
}

func isForBot(msgEvent *slack.MessageEvent) bool {
	callerId := msgEvent.Msg.User
	respondableMessage := msgEvent.Msg.Type == "message" &&
		callerId != botId &&
		msgEvent.Msg.SubType != "message_deleted"
	toBot := strings.Contains(msgEvent.Msg.Text, "<@"+botId+">") ||
		strings.HasPrefix(msgEvent.Msg.Channel, "D")

	return respondableMessage && toBot
}

func main() {
    err := godotenv.Load("phile.env")
    if err != nil {
        log.Fatal("Error loading phile.env file.")
    }
    slackApiToken := os.Getenv("SLACK_API_TOKEN")
    phabApiToken := os.Getenv("PHAB_API_TOKEN")
    phabServerUrl := os.Getenv("PHAB_SERVER_URL")

	slackClient := slack.New(slackApiToken)
	phabClient, err := gonduit.Dial(
		phabServerUrl,
		&core.ClientOptions{APIToken: phabApiToken})
	if err != nil {
        log.Fatal("Error connecting to phabricator conduit.")
	}

	os.Exit(run(slackClient, phabClient))
}
