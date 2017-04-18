package app

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/simon0191/slack-visitor/model"
)

func (app *App) SendChatRequest(visitorName string, subject string) *model.Chat {

	chat := model.Chat{State: "pending", VisitorName: visitorName, Subject: subject}
	if app.db.Create(&chat).Error != nil {
		app.Logger.Printf("Unable to create chat from *%s* with subject _\"%s\"_", visitorName, subject)
		return nil
	}

	text := fmt.Sprintf("New chat request:\n*%s* wants to talk about _\"%s\"_", visitorName, subject)
	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			{
				Text:       "What would you like to do?",
				CallbackID: chat.ID,
				Actions: []slack.AttachmentAction{
					{Name: "chat_request_state", Text: "Accept chat request", Type: "button", Value: "accepted", Style: "primary"},
					{Name: "chat_request_state", Text: "Decline chat request", Type: "button", Value: "declined", Style: "danger"},
				},
			},
		},
	}

	app.SlackBot.PostMessage(app.Config.VisitorChannelID, text, params)
	return &chat
}

func (app *App) GetChatByID(id string) (*model.Chat, error) {
	var chat model.Chat
	if app.db.Where("id = ?", id).First(&chat); app.db.Error != nil {
		return nil, app.db.Error
	}

	return &chat, nil
}
