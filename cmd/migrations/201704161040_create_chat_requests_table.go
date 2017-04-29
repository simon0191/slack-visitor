package migrations

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
	"time"
)

func init() {
	List = append(List, &gormigrate.Migration{
		ID: "201704161040_create_chat_requests_table",
		Migrate: func(tx *gorm.DB) error {

			type Chat struct {
				ID          string  `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
				VisitorName string  `gorm:"type:varchar(100)"`
				Subject     string  `gorm:"type:text"`
				State       string  `gorm:"type:varchar(100)"`
				ChannelID   *string `gorm:"type:varchar(100)"`

				CreatedAt time.Time
				UpdatedAt time.Time
				DeletedAt *time.Time
			}

			if err := tx.AutoMigrate(&Chat{}).Error; err != nil {
				return err
			}

			return tx.Model(Chat{}).AddForeignKey("state", "chat_states (id)", "RESTRICT", "RESTRICT").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("chats").Error
		},
	})
}
