package main

import (
	"fmt"
	"github.com/simon0191/slack-visitor/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
