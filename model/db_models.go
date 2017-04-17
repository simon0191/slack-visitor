package model

import (
	"time"
)

type ChatRequest struct {
	ID        string `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	State     string `gorm:"primary_key;type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
