package app

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/simon0191/slack-visitor/model"
)

func (app *App) SendChatRequest(visitorName string, subject string) *model.ChatRequest {

	request := model.ChatRequest{State: "pending"}
	if app.db.Create(&request).Error != nil {
		app.Logger.Printf("Unable to create chat request from *%s* with subject _\"%s\"_", visitorName, subject)
		return nil
	}

	params := slack.PostMessageParameters{
		Attachments: []slack.Attachment{
			{
				Text:       "What would you like to do?",
				CallbackID: request.ID,
				Actions: []slack.AttachmentAction{
					{Name: "chat", Text: "Accept chat request", Type: "button", Value: "accept", Style: "primary"},
					{Name: "decline", Text: "Decline chat request", Type: "button", Value: "decline", Style: "danger"},
				},
			},
		},
	}

	text := fmt.Sprintf("New chat request:\n*%s* wants to talk about _\"%s\"_", visitorName, subject)
	app.SlackBot.PostMessage(app.Config.VisitorChannelID, text, params)
	return &request
}
