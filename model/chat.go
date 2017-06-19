package model

import (
	"time"
)

const (
	CHAT_STATE_PENDING  = "pending"
	CHAT_STATE_ACCEPTED = "accepted"
	CHAT_STATE_DECLINED = "declined"
	CHAT_STATE_FINISHED = "finished"

	MAX_CHANNEL_NAME_LENGHT = 21
)

type Chat struct {
	ID          string  `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	VisitorName string  `gorm:"type:varchar(100)" json:"visitorName"`
	Subject     string  `gorm:"type:text" json:"subject"`
	State       string  `gorm:"type:varchar(100)" json:"state"`
	ChannelID   *string `gorm:"type:varchar(100)" json:"channelId"`
	MessageTs   *string `gorm:"type:varchar(100)" json:"messageTs"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`

	Messages []Message `json:"messages"`
}
