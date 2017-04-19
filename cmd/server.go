package cmd

import (
	"github.com/simon0191/slack-visitor/api"
	"github.com/simon0191/slack-visitor/app"
	"github.com/simon0191/slack-visitor/model"
	"github.com/simon0191/slack-visitor/utils"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

func runServerCmd(cmd *cobra.Command, args []string) {
	var (
		config *model.Config
		server *api.Server
		a      *app.App
	)

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}
	config, err = utils.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	a = app.New(config)
	server = api.NewServer(config.WebServerSettings, a)
	a.Init()
	server.Run()
}
