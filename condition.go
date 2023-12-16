package main

import (
	"fmt"
	"strconv"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func getCondition(bot *traqwsbot.Bot, channelID string) {
	var conditions []Condition
	if err := db.Select(&conditions, "SELECT * FROM `condition`"); err != nil {
		fmt.Println(err)
	}

	res := "## 状況一覧\n|id|user|状況名|\n|---|---|---|\n"
	var idstr string

	for _, v := range conditions {
		idstr = strconv.Itoa(v.Id)
		res += "|" + idstr + "|" + v.User + "|" + v.Name + "|\n"
	}
	simplePost(bot, channelID, res)

}

func postCondition(bot *traqwsbot.Bot, channelID string, conditionreq string) {
	var conditions []Condition
	if err := db.Select(&conditions, "SELECT * FROM `condition`"); err != nil {
		fmt.Println(err)
	}
	_, err := db.Exec("INSERT INTO `condition` (`user`,`condition`) VALUES(?,?)", channelID, conditionreq)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Failed to add condition")
	}

	//追加した状況のcondition_idを取得
	var condition int
	err = db.Get(&condition, "SELECT `condition_id` FROM `condition` WHERE `condition` = ? ORDER BY `condition_id` DESC", conditionreq)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Failed to add condition")
	}

	conditionstr := strconv.Itoa(condition)
	//simplePost(bot, channelID, "状況の追加が完了しましあ。内容は以下の通りです\n| Condition_id | Condition |\n| --- | --- |\n| "+conditionstr+" | " + conditionreq + " |\n ")
	simplePost(bot, channelID, "状況の追加が完了しましあ。内容は以下の通りです\n| Condition_id | user | Condition |\n| --- | --- | --- |\n| "+conditionstr+" | " + channelID +" | " + conditionreq + " |\n ")

}


