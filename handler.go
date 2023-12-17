package main

import (
	"fmt"
	"strconv"
	"time"

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

func postTask(bot *traqwsbot.Bot, userId string, channelID string, newTask TaskWithoutId) {
	var dateOfNow = time.Now().Format("2006-01-02")
	res, err := db.Exec("INSERT INTO task (user, title, description, condition_id, difficulty, created_at, updated_at, dueDate) VALUES (?,?,?,?,?,?,?,?)", userId, newTask.Title, newTask.Description, newTask.ConditionId, newTask.Difficulty, dateOfNow, dateOfNow, newTask.DueDate)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Datebase insert error")
	}

	taskId, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		simpleEdit(bot, channelID, "ID get error")
	}

	resStr := "## タスクを追加しました！\n|ID|タスク名|期限|追加日|\n|---|---|---|---|\n|" + strconv.Itoa(int(taskId)) + "|" + newTask.Title + "|" + newTask.DueDate + "|" + dateOfNow + "|\n"
	simplePost(bot, channelID, resStr)
}

func deleteTask(bot *traqwsbot.Bot, userID string, channelID string,  taskIDstr string) {
	//消去対象タスクの取得
	taskid, err := strconv.Atoi(taskIDstr)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Please enter a valid taskid")
		return
	}

	var task Task
	err = db.Get(&task, "SELECT * FROM `task` WHERE `id` =? ", taskid)

	//存在判定,他ユーザーのコンディションを消去不可
	if err != nil || task.User != userID {
		if err != nil {
			fmt.Println(err)
		}
		simplePost(bot, channelID, "There is no such task of yours")
		return
	}

	//消去の実行
	_, err = db.Exec("DELETE FROM `task` WHERE `id` =? ", taskid)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Failed to delete task")
		return
	}

	//ログの出力
	simplePost(bot, channelID, "タスクの消去が完了しました。消去内容は以下の通りです\n| ID | タスク名 |\n| --- | --- |\n| "+taskIDstr+" | "+task.Title+" |\n ")
}
