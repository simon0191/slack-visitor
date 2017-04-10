package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/simon0191/slack-visitor/utils"
	"gopkg.in/gormigrate.v1"
)

func main() {
	config, err := utils.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(config.DBSettings.Driver, config.DBSettings.Connection)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	db.LogMode(true)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create chats table
		{
			ID: "201704091117",
			Migrate: func(tx *gorm.DB) error {
				type Chat struct {
					gorm.Model
					UID       string `sql:"index"`
					channel   string
					VisitorID uint `sql:"index"`
				}
				return tx.AutoMigrate(&Chat{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("chats").Error
			},
		},
		// create messages table
		{
			ID: "201704091119",
			Migrate: func(tx *gorm.DB) error {
				type Message struct {
					gorm.Model
					Text     string
					ChatID   uint
					UserType uint
					UserID   string
				}

				if err := tx.AutoMigrate(&Message{}).Error; err != nil {
					return err
				}

				return tx.Model(Message{}).AddForeignKey("chat_id", "chats (id)", "RESTRICT", "RESTRICT").Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("messages").Error
			},
		},
	})

	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Printf("Migration did run successfully")
}
