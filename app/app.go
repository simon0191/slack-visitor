package app

import (
	"github.com/nlopes/slack"
	"github.com/simon0191/slack-visitor/model"

	"github.com/jinzhu/gorm"
	"log"
	"os"
)

type App struct {
	Config   *model.Config
	SlackApp *slack.Client
	SlackBot *slack.Client
	Logger   *log.Logger

	db            *gorm.DB
	botInfo       *slack.Info
	slackRTM      *slack.RTM
	slackChannels map[string]bool

	SendHostMessage      chan *model.Message
	sendVisitorMessage   chan *model.Message
	registerSlackChannel chan string
}

func New(config *model.Config) *App {

	app := &App{
		Config:        config,
		SlackApp:      slack.New(config.SlackSettings.AppAPIKey),
		SlackBot:      slack.New(config.SlackSettings.BotAPIKey),
		Logger:        log.New(os.Stdout, "slack-visitor: ", log.Lshortfile|log.LstdFlags),
		slackChannels: map[string]bool{},

		SendHostMessage:      make(chan *model.Message),
		sendVisitorMessage:   make(chan *model.Message),
		registerSlackChannel: make(chan string),
	}

	slack.SetLogger(app.Logger)
	app.SlackApp.SetDebug(false)
	app.SlackBot.SetDebug(false)

	return app
}

func (app *App) Init() {
	app.InitDB()
	go app.listenRegisterSlackChannel()
	go app.ReadPump()
}

func (app *App) ListenToMessages() {
	/*
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
	*/
}

func (app *App) ReadPump() {
	app.slackRTM = app.SlackBot.NewRTM()
	go app.slackRTM.ManageConnection()

	for msg := range app.slackRTM.IncomingEvents {

		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			//TODO: get the info of Bot user in a different way
			app.botInfo = app.slackRTM.GetInfo()
			app.Logger.Println("Bot Info received")
			app.slackRTM.SendMessage(app.slackRTM.NewOutgoingMessage("Hello world", "C04230HEX"))

		case *slack.MessageEvent:
			app.Logger.Printf("Message: %v\n", ev)
		/*
			if bridge, ok := app.bridges[ev.Channel]; ok {
				bridge.toClient <- []byte(ev.Text)
			}
		*/
		case *slack.PresenceChangeEvent:
			app.Logger.Printf("Presence Change: %v\n", ev)

		case *slack.RTMError:
			app.Logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			app.Logger.Printf("Invalid credentials")
			return
		}
	}
}

func (app *App) InitDB() {
	db, err := gorm.Open(app.Config.DBSettings.Driver, app.Config.DBSettings.Connection)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	app.db = db
	app.Logger.Println("DB initialized")
}

func (app *App) listenRegisterSlackChannel() {
	for channelID := range app.registerSlackChannel {
		app.slackChannels[channelID] = true
	}
}
