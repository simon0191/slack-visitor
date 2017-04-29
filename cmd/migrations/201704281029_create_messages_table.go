package migrations

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"time"
)

func init() {
	List = append(List, &gormigrate.Migration{
		ID: "201704281029_create_messages_table",
		Migrate: func(tx *gorm.DB) error {
			type Message struct {
				ID       string `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
				ChatID   string `gorm:"type:uuid"`
				Content  string `gorm:"type:varchar(100)"`
				Source   string `gorm:"type:varchar(100)"`
				FromName string `gorm:"type:varchar(100)"`

				CreatedAt time.Time
				UpdatedAt time.Time
				DeletedAt *time.Time
			}

			if err := tx.AutoMigrate(&Message{}).Error; err != nil {
				return err
			}

			if err := tx.Model(Message{}).AddForeignKey("chat_id", "chats (id)", "RESTRICT", "RESTRICT").Error; err != nil {
				return err
			}

			return tx.Model(Message{}).AddForeignKey("source", "message_sources (id)", "RESTRICT", "RESTRICT").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("messages").Error
		},
	})
}
