package main

import (
	"github.com/simon0191/slack-visitor/app"
	"github.com/simon0191/slack-visitor/utils"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	config, err := utils.LoadConfig("./config/config.json")

	if err != nil {
		panic(err)
	}

	a := app.New(config)
	a.Init()

}
