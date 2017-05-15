package app

import (
	"github.com/nlopes/slack"
	"github.com/simon0191/slack-visitor/model"

	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

const (
	REGISTER_IGNORED_MESSAGE_WAIT = 100
)

type App struct {
	Config   *model.Config
	SlackApp *slack.Client
	SlackBot *slack.Client
	Logger   *log.Logger

	db       *gorm.DB
	botInfo  *slack.Info
	slackRTM *slack.RTM

	messagesChs         map[string]chan *model.Message
	messagesSubscribers map[string][]MessageSubscriber

	registerChatCh     chan *model.Chat
	newChatSubscribers []NewChatSubscriber

	ignoreSlackMessages map[string]bool
}

type MessageSubscriber func(m *model.Message)
type NewChatSubscriber func(c *model.Chat)

func New(config *model.Config) *App {

	app := &App{
		Config:   config,
		SlackApp: slack.New(config.SlackSettings.AppAPIKey),
		SlackBot: slack.New(config.SlackSettings.BotAPIKey),
		Logger:   log.New(os.Stdout, "slack-visitor: ", log.Lshortfile|log.LstdFlags),

		messagesChs:         map[string]chan *model.Message{},
		messagesSubscribers: map[string][]MessageSubscriber{},

		registerChatCh:     make(chan *model.Chat),
		newChatSubscribers: []NewChatSubscriber{},

		ignoreSlackMessages: map[string]bool{},
	}

	slack.SetLogger(app.Logger)
	app.SlackApp.SetDebug(false)
	app.SlackBot.SetDebug(false)

	return app
}

func (app *App) Init() {
	app.initDB()
	go app.listenToRegisterChatCh()
	go app.listenToSlackEvents()
}

func (app *App) listenToSlackEvents() {
	app.slackRTM = app.SlackBot.NewRTM()
	go app.slackRTM.ManageConnection()

	for msg := range app.slackRTM.IncomingEvents {

		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			//TODO: get the info of Bot user in a different way
			app.botInfo = app.slackRTM.GetInfo()
			app.Logger.Printf("Bot Info received: %+v\n", app.botInfo.User)

		case *slack.MessageEvent:
			//TODO: better way to avoid repeated messages
			if ev.SubType != "group_archive" {
				go func() {
					time.Sleep(REGISTER_IGNORED_MESSAGE_WAIT * time.Millisecond)
					if _, ok := app.ignoreSlackMessages[ev.Timestamp]; ok {
						delete(app.ignoreSlackMessages, ev.Timestamp)
					} else {
						app.Logger.Printf("Message: %+v\n", ev)
						go app.SendHostMessage(ev)
					}
				}()
			}

		case *slack.GroupArchiveEvent:
			chat, err := app.GetChatByChannel(ev.Channel)
			if err == nil {
				app.TerminateChat(chat, false)
			}

		case *slack.RTMError:
			app.Logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			app.Logger.Printf("Invalid credentials")
			return
		}
	}
}

func (app *App) initDB() {
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

func (app *App) listenToRegisterChatCh() {
	for chat := range app.registerChatCh {
		app.messagesChs[chat.ID] = make(chan *model.Message)
		app.messagesSubscribers[chat.ID] = []MessageSubscriber{}

		app.OnMessage(chat.ID, func(message *model.Message) {
			if message.Source == model.MESSAGE_SOURCE_VISITOR {
				_, timestamp, err := app.SlackBot.PostMessage(*message.Chat.ChannelID, message.Content, slack.PostMessageParameters{
					Username: message.Chat.VisitorName,
					AsUser:   false,
				})

				if err != nil {
					app.Logger.Printf("Unable to send message: %s\n", err)
					return
				}
				app.ignoreSlackMessages[timestamp] = true
				//Clean ignoreSlackMessages
				go func() {
					time.Sleep(REGISTER_IGNORED_MESSAGE_WAIT * time.Millisecond * 2)
					delete(app.ignoreSlackMessages, timestamp)
				}()
			}
		})

		go app.listenToMessages(chat.ID)
		for _, notify := range app.newChatSubscribers {
			notify(chat)
		}
	}
}

func (app *App) listenToMessages(chatID string) {
	for message := range app.messagesChs[chatID] {
		for _, notify := range app.messagesSubscribers[chatID] {
			go notify(message)
		}
	}
}

func (app *App) OnMessage(chatID string, subscriber MessageSubscriber) {
	app.messagesSubscribers[chatID] = append(app.messagesSubscribers[chatID], subscriber)
}

func (app *App) OnNewChat(subscriber NewChatSubscriber) {
	app.newChatSubscribers = append(app.newChatSubscribers, subscriber)
}
