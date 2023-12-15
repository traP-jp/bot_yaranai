package main

import (
	"fmt"
	"os"
	"strings"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}
	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		fmt.Println(p.Message.Text)
		cmd := strings.Split(p.Message.Text, " ")
		if cmd[1] == "stamp" {
			if cmd[2] == "recent" {
				if len(cmd) > 3 && cmd[3] == "user" {
					checkUserHandrer(bot, p)
				} else {
					checkHandrer(bot, p)
				}
			} else if cmd[2] == "count" {
				stampCountHandrer(bot, p)
			} else {
				simplePost(bot, p.Message.ChannelID, "No such command")
			}
		} else if cmd[1] == "heatmap" {
			heatMapHandrer(bot, p)
		} else if cmd[1] == "help" {
			bytes, err := os.ReadFile("help.txt")
			if err != nil {
				panic(err)
			}
			simplePost(bot, p.Message.ChannelID, string(bytes))
		} else {
			simplePost(bot, p.Message.ChannelID, "No such command")
		}

	})

	if err := bot.Start(); err != nil {
		panic(err)
	}
}
