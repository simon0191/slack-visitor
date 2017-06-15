package cmd

import (
	"github.com/simon0191/slack-visitor/api"
	"github.com/simon0191/slack-visitor/app"
	"github.com/simon0191/slack-visitor/config"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

func runServerCmd(cmd *cobra.Command, args []string) {
	var (
		c      *config.Config
		server *api.Server
		a      *app.App
	)

	c = config.Load()

	rand.Seed(time.Now().UnixNano())

	a = app.New(c)
	server = api.NewServer(c.WebServerSettings, a)
	a.Init()
	server.Run()
}
