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

func postTask(bot *traqwsbot.Bot, channelID string, newTask TaskWithoutId) {
	var dateOfNow = time.Now().Format("2006-01-02")
	var userId = "testuser"
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

	resStr := "##タスクを追加しました！\n|ID|タスク名|期限|\n|---|---|---|\n|" + strconv.Itoa(int(taskId)) + "|" + newTask.Title + "|" + newTask.DueDate + "|\n"
	simplePost(bot, channelID, resStr)
}
