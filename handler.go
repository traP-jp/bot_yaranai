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

func getTask(bot *traqwsbot.Bot, userID string, channelID string) {
	var tasks []Task
	if err := db.Select(&tasks, "SELECT * FROM `task` WHERE user = ?", userID); err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "There is no such task of yours")
		return
	}

	res := "## タスク一覧\n|ID|タスク名|期限|\n|---|---|---|\n"
	for _, v := range tasks {
		idStr := strconv.Itoa(v.Id)
		dueDateStr := v.DueDate.Format("2006-01-02")
		res += "|" + idStr + "|" + v.Title + "|" + dueDateStr + "|\n"
	}
	simplePost(bot, channelID, res)
}

func postTask(bot *traqwsbot.Bot, userID string, channelID string, newTask TaskWithoutId) {
	var dateOfNow = time.Now().Format("2006-01-02")
	res, err := db.Exec("INSERT INTO task (user, title, description, condition_id, difficulty, created_at, updated_at, dueDate) VALUES (?,?,?,?,?,?,?,?)", userID, newTask.Title, newTask.Description, newTask.ConditionId, newTask.Difficulty, dateOfNow, dateOfNow, newTask.DueDate)

	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Datebase insert error")
		return
	}

	taskID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		simpleEdit(bot, channelID, "ID get error")
		return
	}

	resStr := "## タスクを追加しました！\n|ID|タスク名|期限|\n|---|---|---|\n|" + strconv.Itoa(int(taskID)) + "|" + newTask.Title + "|" + newTask.DueDate + "|\n"
	simplePost(bot, channelID, resStr)
}

func putTask(bot *traqwsbot.Bot, taskID int, userID string, channelID string, changeList [][2]string) {
	var existingTask Task
	err := db.Get(&existingTask, "SELECT * FROM task WHERE `id` = ? AND `user` = ?", taskID, userID)
	if err != nil {
		fmt.Println(taskID)
		fmt.Println(userID)
		fmt.Println(err)
		simplePost(bot, channelID, "無効なタスク ID です")
		return
	}

	for _, query := range changeList {
		var err error
		if query[0] == "condition_id" || query[0] == "difficulty" {
			queryInt, err := strconv.Atoi(query[1])
			if err != nil && query[0] == "condition_id" {
				fmt.Println(query[1])
				fmt.Println(err)
				simplePost(bot, channelID, "無効な Condition ID です。")
				return
			} else if err != nil {
				fmt.Println(err)
				simplePost(bot, channelID, "無効な Difficluty です。")
				return
			}
			_, err = db.Exec("UPDATE task SET ? = ? WHERE id = ? AND user = ?", query[0], queryInt, taskID, userID)
			if err != nil {
				fmt.Println(query[0])
				fmt.Println(queryInt)
				fmt.Println(err)
				simplePost(bot, channelID, "Internal Server Error")
				return
			}
		} else {
			_, err = db.Exec("UPDATE task SET ? = ? WHERE id = ? AND user = ?", query[0], query[1], taskID, userID)
			if err != nil {
				fmt.Println(query[0])
				fmt.Println(query[1])
				fmt.Println(err)
				simplePost(bot, channelID, "Internal Server Error")
				return
			}
		}
		nowOfDate := time.Now()
		_, err = db.Exec("UPDATE task SET updated_at = ? WHERE id = ? AND user = ?", nowOfDate, taskID, userID)
		if err != nil {
			fmt.Println(err)
			simplePost(bot, channelID, "Internal Server Error")
			return
		}
	}

	conditionIdInt := strconv.Itoa(existingTask.ConditionId)
	difficultyInt := strconv.Itoa(existingTask.Difficulty)

	dueDateStr := existingTask.DueDate.Format("2006-01-02")
	resStr := "## タスクを変更しました。\n変更結果\n|タスク名|詳細|状況|こなしにくさ|期限|\n|---|---|---|---|---|\n|" + existingTask.Title + "|" + existingTask.Description + "|" + conditionIdInt + "|" + difficultyInt + "|" + dueDateStr + "|\n"

	simplePost(bot, channelID, resStr)
}
