package app

import (
	"github.com/simon0191/slack-visitor/model"
)

func (app *App) SendVisitorMessage(chat *model.Chat, msg string) {
	message := &model.Message{
		ChatID:  chat.ID,
		Content: msg,
		Source:  model.MESSAGE_SOURCE_VISITOR,
	}
	app.db.Create(message)
	app.sendVisitorMessage <- message
}
