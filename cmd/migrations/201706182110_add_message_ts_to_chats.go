package migrations

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func init() {
	List = append(List, &gormigrate.Migration{
		ID: "201706182110_add_message_ts_to_chats",
		Migrate: func(tx *gorm.DB) error {
			type Chat struct {
				MessageTs *string `gorm:"type:varchar(100)" json:"messageTs"`
			}

			return tx.AutoMigrate(&Chat{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Table("chats").DropColumn("message_ts").Error
		},
	})
}
