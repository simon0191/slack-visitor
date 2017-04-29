package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/simon0191/slack-visitor/model"
	"gopkg.in/gormigrate.v1"
)

func init() {
	List = append(List, &gormigrate.Migration{
		ID: "201704281028_create_message_sources_table",
		Migrate: func(tx *gorm.DB) error {
			type MessageSource struct {
				ID string `gorm:"primary_key;type:varchar(100)"`
			}
			if err := tx.AutoMigrate(&MessageSource{}).Error; err != nil {
				return err
			}

			sources := []string{
				model.MESSAGE_SOURCE_SLACK,
				model.MESSAGE_SOURCE_VISITOR,
			}

			for _, source := range sources {
				if err := tx.Create(&MessageSource{ID: source}).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("message_sources").Error
		},
	})
}
