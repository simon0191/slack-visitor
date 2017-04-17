package model

import (
	"github.com/jinzhu/gorm"
)

type Chat struct {
	gorm.Model
	UID       string `sql:"index"`
	VisitorID uint   `sql:"index"`
}

type Message struct {
	gorm.Model
	Text     string
	Chat     Chat `gorm:"ForeignKey:ChatID"`
	ChatID   uint
	UserType string
	UserID   string
}
