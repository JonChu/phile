package main

import (
    "log"
    "os"

    "github.com/nlopes/slack"
)

const API_TOKEN = "<INSERT API TOKEN HERE>"

var (
        api *slack.Client
        botId string
)

func run(api *slack.Client) int {
    log.Printf("Setting up RTM connection...")

    rtm := api.NewRTM()
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
                rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", ev.Channel))
            
            case *slack.RTMError:
                log.Printf("Error: %s\n", ev.Error())
            
            case *slack.InvalidAuthEvent:
                log.Print("Invalid credentials")
                return 1
            }
        }
    }
}

func main() {
    log.SetOutput(os.Stdout)
    api := slack.New(API_TOKEN)
    os.Exit(run(api))
}