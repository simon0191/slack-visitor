package app

import (
	"github.com/nlopes/slack"
	"github.com/simon0191/slack-visitor/model"
	"time"
)

func (app *App) SendVisitorMessage(chatID, msg string) {
	chat, err := app.GetChatByID(chatID)
	if err != nil {
		app.Logger.Printf("Invalid Chat id: %s\n", chatID)
		return
	}

	message := &model.Message{
		ChatID:  chat.ID,
		Content: msg,
		Source:  model.MESSAGE_SOURCE_VISITOR,
	}
	if err := app.db.Create(message).Error; err != nil {
		app.Logger.Printf("Unable to create message: %+v\n%s\n", message, err)
		return
	}
	message.Chat = *chat

	app.sendMessage(message)
}

func (app *App) SendHostMessage(slackMsg *slack.MessageEvent) {
	chat, err := app.GetChatByChannel(slackMsg.Channel)
	if err != nil {
		app.Logger.Print("Invalid Slack Channel " + slackMsg.Channel)
		return
	}

	//TODO: retrieve user info from cache
	user, err := app.SlackApp.GetUserInfo(slackMsg.User)
	if err != nil {
		app.Logger.Printf("Unable to retrieve user info: %s\n", err)
		return
	}

	message := &model.Message{
		ChatID:   chat.ID,
		Content:  slackMsg.Text,
		Source:   model.MESSAGE_SOURCE_SLACK,
		FromName: user.RealName,
	}

	if user.RealName == "" {
		message.FromName = user.Name
	}

	if err := app.db.Create(message).Error; err != nil {
		app.Logger.Printf("Unable to create message: %+v\n%s\n", message, err)
		return
	}
	message.Chat = *chat

	app.sendMessage(message)
}

func (app *App) sendMessage(message *model.Message) {
	ch, ok := app.messagesChs[message.ChatID]
	if !ok {
		app.registerChat(&message.Chat)
		// Wait until chat channel is created
		for {
			time.Sleep(100 * time.Millisecond)
			if ch, ok = app.messagesChs[message.ChatID]; ok {
				break
			}
		}
	}

	ch <- message
}
