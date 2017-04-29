package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/simon0191/slack-visitor/model"
	"gopkg.in/gormigrate.v1"
)

func init() {
	List = append(List, &gormigrate.Migration{
		ID: "201704161038_create_chat_states_table",
		Migrate: func(tx *gorm.DB) error {
			type ChatState struct {
				ID string `gorm:"primary_key;type:varchar(100)"`
			}
			if err := tx.AutoMigrate(&ChatState{}).Error; err != nil {
				return err
			}
			for _, state := range []string{model.CHAT_STATE_ACCEPTED, model.CHAT_STATE_DECLINED, model.CHAT_STATE_FINISHED, model.CHAT_STATE_PENDING} {
				if err := tx.Create(&ChatState{ID: state}).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("chat_states").Error
		},
	})
}
