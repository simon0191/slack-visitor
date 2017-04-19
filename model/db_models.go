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
	VisitorName string  `gorm:"type:varchar(100)" json:"visitor_name"`
	Subject     string  `gorm:"type:text" json:"subject"`
	State       string  `gorm:"type:varchar(100)" json:"state"`
	ChannelID   *string `gorm:"type:varchar(100)" json:"channel_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
