package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "config/config.json", "Configuration file to use.")

	RootCmd.AddCommand(
		dbCmd,
	)
}

var RootCmd = &cobra.Command{
	Use:   "visitor",
	Short: "Open source, self-hosted Slack visitor chat",
	Run:   runServerCmd,
}
