package app

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/simon0191/slack-visitor/model"
	"github.com/simon0191/slack-visitor/utils"
	"regexp"
	"strings"
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
					{Name: "action", Text: "Accept chat request", Type: "button", Value: "accept_chat", Style: "primary"},
					{Name: "action", Text: "Decline chat request", Type: "button", Value: "decline_chat", Style: "danger"},
				},
			},
		},
	}

	app.SlackBot.PostMessage(app.Config.VisitorChannelID, text, params)
	return &chat
}

func (app *App) GetChatByID(id string) (*model.Chat, error) {
	var chat model.Chat
	if err := app.db.Where("id = ?", id).First(&chat).Error; err != nil {
		return nil, err
	}

	return &chat, nil
}

func (app *App) AcceptChat(action slack.AttachmentActionCallback) {
	var chat model.Chat
	if err := app.db.Where("id = ? AND state = ?", action.CallbackID, model.CHAT_STATE_PENDING).First(&chat).Error; err != nil {
		app.Logger.Println(err)
		return
	}

	group := app.createSlackGroup(&chat)
	app.SlackApp.InviteUserToGroup(group.ID, app.botInfo.User.ID)
	app.SlackApp.InviteUserToGroup(group.ID, action.User.ID)
	app.updateAcceptedMessage(action.Channel.ID, action.User, action.MessageTs, chat)

	chat.State = model.CHAT_STATE_ACCEPTED
	chat.ChannelID = &group.ID

	if err := app.db.Save(&chat).Error; err != nil {
		app.Logger.Println(err)
	}
}

func (app *App) DeclineChat(action slack.AttachmentActionCallback) {
	var chat model.Chat
	if err := app.db.Where("id = ? AND state = ?", action.CallbackID, model.CHAT_STATE_PENDING).First(&chat).Error; err != nil {
		app.Logger.Println(err)
		return
	}

	app.updateDeclinedMessage(action.Channel.ID, action.User, action.MessageTs, chat)
	chat.State = model.CHAT_STATE_DECLINED
	if err := app.db.Save(&chat).Error; err != nil {
		app.Logger.Println(err)
	}
}

func (app *App) JoinChat(action slack.AttachmentActionCallback) {
	var chat model.Chat
	if err := app.db.Where("id = ? AND state = ?", action.CallbackID, model.CHAT_STATE_ACCEPTED).First(&chat).Error; err != nil {
		app.Logger.Println(err)
		return
	}

	app.SlackApp.InviteUserToGroup(*chat.ChannelID, action.User.ID)
}

func (app *App) updateAcceptedMessage(channelID string, user slack.User, messageTs string, chat model.Chat) {
	app.SlackBot.SendMessage(
		channelID,
		slack.MsgOptionText(fmt.Sprintf("Accepted chat request:\n*%s* wants to talk about _\"%s\"_", chat.VisitorName, chat.Subject), false),
		slack.MsgOptionUpdate(messageTs),
		slack.MsgOptionAttachments(slack.Attachment{
			Text:       fmt.Sprintf("@%s has accepted this chat request", user.Name),
			CallbackID: chat.ID,
			Actions: []slack.AttachmentAction{
				{Name: "action", Text: "Join this chat", Type: "button", Value: "join_chat", Style: "primary"},
			},
		}),
	)
}

func (app *App) updateDeclinedMessage(channelID string, user slack.User, messageTs string, chat model.Chat) {

	app.SlackBot.SendMessage(
		channelID,
		slack.MsgOptionText(fmt.Sprintf("Declined chat request:\n*%s* wants to talk about _\"%s\"_", chat.VisitorName, chat.Subject), false),
		slack.MsgOptionUpdate(messageTs),
		slack.MsgOptionAttachments(slack.Attachment{
			Text:       fmt.Sprintf("@%s has declined this chat request", user.Name),
			CallbackID: chat.ID,
		}),
	)
}

func (app *App) createSlackGroup(chat *model.Chat) *slack.Group {
	group, err := app.SlackApp.CreateGroup(buildChannelName(chat.VisitorName))
	if err != nil {
		app.Logger.Fatal(err)
	}

	_, err = app.SlackApp.SetGroupPurpose(group.ID, chat.Subject)
	if err != nil {
		app.Logger.Fatal(err)
	}

	_, err = app.SlackApp.SetGroupTopic(group.ID, chat.Subject)
	if err != nil {
		app.Logger.Fatal(err)
	}

	app.registerSlackChannel <- group.ID

	return group
}

func buildChannelName(visitorName string) string {
	regex := regexp.MustCompile("\\W")
	visitorName = strings.ToLower(visitorName)
	channelName := utils.RandString(4) + "-" + regex.ReplaceAllLiteralString(visitorName, "-")

	if len(channelName) > model.MAX_CHANNEL_NAME_LENGHT {
		channelName = channelName[0:model.MAX_CHANNEL_NAME_LENGHT]
	}

	return channelName
}
