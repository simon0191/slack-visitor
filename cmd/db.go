package cmd

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/simon0191/slack-visitor/utils"
	"github.com/spf13/cobra"
	"gopkg.in/gormigrate.v1"
	"time"
)

func init() {
	dbCmd.AddCommand(
		dbMigrateCmd,
		dbRollbackCmd,
	)
}

var dbCmd = &cobra.Command{
	Use: "db",
}

var dbMigrateCmd = &cobra.Command{
	Use: "migrate",
	Run: dbMigrateCmdFunc,
}

var dbRollbackCmd = &cobra.Command{
	Use: "rollback",
	Run: dbRollbackCmdFunc,
}

//TODO: place each migration in a different file
var migrations = []*gormigrate.Migration{
	// load pgcrypto extension
	{
		ID: "201704160930",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("CREATE EXTENSION pgcrypto;").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("DROP EXTENSION IF EXISTS pgcrypto;").Error
		},
	},
	// create chat request states table
	{
		ID: "201704161038",
		Migrate: func(tx *gorm.DB) error {
			type ChatState struct {
				ID string `gorm:"primary_key;type:varchar(100)"`
			}
			if err := tx.AutoMigrate(&ChatState{}).Error; err != nil {
				return err
			}
			for _, state := range []string{"pending", "accepted", "declined"} {
				if err := tx.Create(&ChatState{ID: state}).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("chat_states").Error
		},
	},
	// create chat requests table
	{
		ID: "201704161040",
		Migrate: func(tx *gorm.DB) error {

			type Chat struct {
				ID          string `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
				VisitorName string `gorm:"type:varchar(100)"`
				Subject     string `gorm:"type:text"`
				State       string `gorm:"primary_key;type:varchar(100)"`

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
	},
}

func dbMigrateCmdFunc(cmd *cobra.Command, args []string) {
	m := initMigrations(cmd)
	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	log.Printf("Migration did run successfully")
}

func dbRollbackCmdFunc(cmd *cobra.Command, args []string) {
	//TODO: implement Rollback specifying migration with m.RollbackMigration(migration)
	m := initMigrations(cmd)
	if err := m.RollbackLast(); err != nil {
		log.Fatalf("Could not rollback: %v", err)
	}

	log.Printf("Rollback did run successfully")
}

func initMigrations(cmd *cobra.Command) *gormigrate.Gormigrate {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatal(err)
	}

	config, err := utils.LoadConfig(configPath)
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

	return gormigrate.New(db, gormigrate.DefaultOptions, migrations)
}
