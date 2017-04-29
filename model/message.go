package model

import (
	"time"
)

const (
	MESSAGE_SOURCE_SLACK   = "slack"
	MESSAGE_SOURCE_VISITOR = "visitor"
)

type Message struct {
	ID       string `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	ChatID   string `gorm:"type:varchar(100)" json:"chatId"`
	Source   string `gorm:"type:varchar(100)" json:"source"`
	FromName string `gorm:"type:varchar(100)" json:"fromName"`
	Content  string `gorm:"type:text" json:"content"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`

	Chat Chat
}
