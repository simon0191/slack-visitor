package model

import (
	"time"
)

type Chat struct {
	ID          string `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	VisitorName string `gorm:"type:varchar(100)" json:"visitor_name"`
	Subject     string `gorm:"type:text" json:"subject"`
	State       string `gorm:"primary_key;type:varchar(100)" json:"state"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
