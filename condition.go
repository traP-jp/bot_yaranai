package main

import (
	"fmt"
	"strconv"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func getCondition(bot *traqwsbot.Bot, channelID string, userID string) {
	//状況リストの取得 User単位(traQチャンネルのUUID)で
	var conditions []Condition
	if err := db.Select(&conditions, "SELECT * FROM `condition` WHERE `user`=?", userID); err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "There is no such condition of yours")
		return
	}

	//状況リストの生成
	res := "## 状況一覧\n|id|状況名|\n|---|---|\n"
	var idstr string

	for _, v := range conditions {
		idstr = strconv.Itoa(v.Id)
		res += "|" + idstr + "|" + v.Name + "|\n"
	}
	simplePost(bot, channelID, res)

}

func postCondition(bot *traqwsbot.Bot, channelID string, conditionreq string,userID string) {

	//引数に持った情報をもとに新規作成
	_, err := db.Exec("INSERT INTO `condition` (`user`,`condition`) VALUES(?,?)", userID, conditionreq)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Failed to add condition")
		return
	}

	//追加した状況のcondition_idを取得
	var condition int
	err = db.Get(&condition, "SELECT `condition_id` FROM `condition` WHERE `condition` = ? ORDER BY `condition_id` DESC", conditionreq)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Failed to add condition")
		return
	}

	//ログの出力
	conditionstr := strconv.Itoa(condition)
	simplePost(bot, channelID, "状況の追加が完了しました。内容は以下の通りです\n| Condition_id | Condition |\n| --- | --- |\n| "+conditionstr+" | "+conditionreq+" |\n ")

}

func putCondition(bot *traqwsbot.Bot, channelID string, conditionidstr string, conditionreq string,UserID string) {
	//更新対象コンディションの取得
	conditionid, err := strconv.Atoi(conditionidstr)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Please enter a valid condition_id\n`condition edit {id} {hoge}`")
		return
	}

	var condition Condition
	err = db.Get(&condition, "SELECT * FROM `condition` WHERE `condition_id` =? ", conditionid)

	//ユーザ同一性の確認
	if err != nil || condition.User != UserID {
		if err != nil {
			fmt.Println(err)
		}
		simplePost(bot, channelID, "There is no such condition of yours")
		return
	}

	//UPDATEの実行
	_, err = db.Exec("UPDATE `condition` set `condition`=? WHERE `condition_id` =?", conditionreq, conditionid)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Failed to rename condition")
		return
	}

	simplePost(bot, channelID, "状況の編集が完了しました。内容は以下の通りです\n| Condition_id | Condition |\n| --- | --- |\n| "+conditionidstr+" | "+conditionreq+" |\n ")

}

func deleteCondition(bot *traqwsbot.Bot, channelID string, conditionidstr string, UserID string) {
	//消去対象コンディションの取得
	conditionid, err := strconv.Atoi(conditionidstr)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Please enter a valid condition_id")
		return
	}

	var condition Condition
	err = db.Get(&condition, "SELECT * FROM `condition` WHERE `condition_id` =? ", conditionid)

	//存在判定,他ユーザーのコンディションを消去不可
	if err != nil || condition.User != UserID {
		if err != nil {
			fmt.Println(err)
		}
		simplePost(bot, channelID, "There is no such condition of yours")
		return
	}

	//消去の実行
	_, err = db.Exec("DELETE FROM `condition` WHERE `condition_id` =? ", conditionid)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Failed to delete condition")
		return
	}

	//ログの出力
	simplePost(bot, channelID, "状況の消去が完了しました。消去内容は以下の通りです\n| Condition_id | Condition |\n| --- | --- |\n| "+conditionidstr+" | "+condition.Name+" |\n ")
}

// デバッグ用全コンディション取得関数
func debuggetCondition(bot *traqwsbot.Bot, channelID string) {
	//状況リストの取得
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
