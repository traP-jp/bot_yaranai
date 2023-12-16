package main

import (
	"fmt"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func getTest(bot *traqwsbot.Bot, channelID string) {
	var tasks []Task
	if err := db.Select(&tasks, "SELECT * FROM task"); err != nil {
		fmt.Println(err)
	}
	res := "## タスク一覧\n|タスク名|期限|\n|---|---|\n"
	for _, v := range tasks {
		res += "|" + v.Title + "|" + v.Description + "|\n"
	}
	simplePost(bot, channelID, res)

}


