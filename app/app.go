package app

import (
	"github.com/nlopes/slack"
	"github.com/simon0191/slack-visitor/model"

	"log"
	"os"
)

type ClientMessage struct {
	channel string
	message string
}

type App struct {
	Config   *model.Config
	SlackApp *slack.Client
	SlackBot *slack.Client
	Logger   *log.Logger

	bridges          map[string]*Bridge
	registerClient   chan *Bridge
	unregisterClient chan *Bridge
	toSlack          chan *ClientMessage
	botInfo          *slack.Info
	slackRTM         *slack.RTM
}

func New(config *model.Config) *App {

	app := &App{
		Config:           config,
		SlackApp:         slack.New(config.SlackSettings.AppAPIKey),
		SlackBot:         slack.New(config.SlackSettings.BotAPIKey),
		Logger:           log.New(os.Stdout, "slack-visitor: ", log.Lshortfile|log.LstdFlags),
		bridges:          map[string]*Bridge{},
		registerClient:   make(chan *Bridge),
		unregisterClient: make(chan *Bridge),
		toSlack:          make(chan *ClientMessage),
	}

	slack.SetLogger(app.Logger)
	app.SlackApp.SetDebug(config.DebugEnabled)
	app.SlackBot.SetDebug(config.DebugEnabled)

	return app
}

func (app *App) ListenToBridges() {
	for {
		select {
		case bridge := <-app.registerClient:

			channel, err := app.SlackApp.CreateChannel(bridge.channel)
			if err != nil {
				app.Logger.Fatal(err)
				break
			}
			//TODO: avoid modification of bridge
			bridge.channel = channel.ID
			app.bridges[channel.ID] = bridge
			app.SlackApp.InviteUserToChannel(channel.ID, app.botInfo.User.ID)

		case bridge := <-app.unregisterClient:
			if _, ok := app.bridges[bridge.channel]; ok {
				delete(app.bridges, bridge.channel)
				close(bridge.toClient)
			}

		case message := <-app.toSlack:

			app.slackRTM.SendMessage(app.slackRTM.NewOutgoingMessage(message.message, message.channel))
		}
	}
}

func (app *App) ReadPump() {
	app.slackRTM = app.SlackBot.NewRTM()
	go app.slackRTM.ManageConnection()

	for msg := range app.slackRTM.IncomingEvents {

		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:

			app.botInfo = app.slackRTM.GetInfo()
			app.slackRTM.SendMessage(app.slackRTM.NewOutgoingMessage("Hello world", "C04230HEX"))

		case *slack.MessageEvent:
			app.Logger.Printf("Message: %v\n", ev)
			if bridge, ok := app.bridges[ev.Channel]; ok {
				bridge.toClient <- []byte(ev.Text)
			}

		case *slack.PresenceChangeEvent:
			app.Logger.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			app.Logger.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			app.Logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			app.Logger.Printf("Invalid credentials")
			return
		}
	}
}
