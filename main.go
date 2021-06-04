package main

import (
	"fmt"

	discordbot "github.com/ab185508/weather-discord-bot/subfiles"
)

func main() {
	err := discordbot.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	discordbot.Start()

	<-make(chan struct{})
}
