package main

import (
	"log"
	"os"
	"strings"

	"github.com/etcinit/gonduit"
	"github.com/etcinit/gonduit/core"
	"github.com/nlopes/slack"
)

const SLACK_API_TOKEN = "xoxb-152541874545-1RHxG8fucCGINtLxZrYNV7KQ"
const PHAB_API_TOKEN = "api-eaiyktsjdkseg33s36nxug7g3j4t"
const PHAB_INSTALL_URL = "https://opendoor.phacility.com"

var botId string

func run(slackClient *slack.Client, phabClient *gonduit.Conn) int {
	log.Printf("Setting up RTM connection...")

	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

	//log.Printf("Connection successfully established.")
	log.Printf("Accepting incoming events...")

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			//log.Printf("Incoming event received...\n")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID
				log.Printf("Connection established...")
				log.Printf("BotId: %v\n", botId)

			case *slack.HelloEvent:
				log.Printf("Server ack received...")

			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev)

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
				log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
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
	log.SetOutput(os.Stdout)

	slackClient := slack.New(SLACK_API_TOKEN)
	phabClient, err := gonduit.Dial(
		PHAB_INSTALL_URL,
		&core.ClientOptions{APIToken: PHAB_API_TOKEN})
	if err != nil {
		os.Exit(1)
	}

	os.Exit(run(slackClient, phabClient))
}
