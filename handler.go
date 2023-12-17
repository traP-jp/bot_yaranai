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

func putTask(bot *traqwsbot.Bot, taskID int, userID string, channelID string, changeList [5]string) {
	var existingTask Task
	err := db.Get(&existingTask, "SELECT * FROM task WHERE `id` = ? AND `user` = ?", taskID, userID)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "無効なタスク ID です")
		return
	}

	information := [...]string{"title", "description", "condition_id", "difficulty", "dueDate"}
	var check = false
	for i, change := range changeList {
		if change == "_" {
			continue
		} else {
			check = true
			if information[i] == "title" {
				_, err = db.Exec("UPDATE task SET title = ? WHERE id = ? AND user = ?", change, taskID, userID)
				if err != nil {
					fmt.Println(err)
					simplePost(bot, channelID, "Internal Server Error: Title")
					return
				}
			} else if information[i] == "description" {
				_, err = db.Exec("UPDATE task SET description = ? WHERE id = ? AND user = ?", change, taskID, userID)
				if err != nil {
					fmt.Println(err)
					simplePost(bot, channelID, "Internal Server Error: Description")
					return
				}
			} else if information[i] == "condition_id" {
				_, err = db.Exec("UPDATE task SET condition_id = ? WHERE id = ? AND user = ?", change, taskID, userID)
				if err != nil {
					fmt.Println(err)
					simplePost(bot, channelID, "Internal Server Error: Condition ID")
					return
				}
			} else if information[i] == "difficulty" {
				_, err = db.Exec("UPDATE task SET difficulty = ? WHERE id = ? AND user = ?", change, taskID, userID)
				if err != nil {
					fmt.Println(err)
					simplePost(bot, channelID, "Internal Server Error: Difficulty")
					return
				}
			} else if information[i] == "dueDate" {
				_, err = db.Exec("UPDATE task SET dueDate = ? WHERE id = ? AND user = ?", change, taskID, userID)
				if err != nil {
					fmt.Println(err)
					simplePost(bot, channelID, "Internal Server Error: DueDate")
					return
				}
			}
		}
	}

	if !check {
		simplePost(bot, channelID, "変更内容を入力してください。")
		return
	}

	nowOfDate := time.Now()
	_, err = db.Exec("UPDATE task SET updated_at = ? WHERE id = ? AND user = ?", nowOfDate, taskID, userID)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Internal Server Error")
		return
	}

	var updatedTask Task
	err = db.Get(&updatedTask, "SELECT * FROM task WHERE `id` = ? AND `user` = ?", taskID, userID)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "無効なタスク ID です")
		return
	}

	conditionIdInt := strconv.Itoa(updatedTask.ConditionId)
	difficultyInt := strconv.Itoa(updatedTask.Difficulty)

	dueDateStr := updatedTask.DueDate.Format("2006-01-02")

	resStr := "## タスクを変更しました。\n変更結果\n|タスク名|詳細|状況|こなしにくさ|期限|\n|---|---|---|---|---|\n|" + updatedTask.Title + "|" + updatedTask.Description + "|" + conditionIdInt + "|" + difficultyInt + "|" + dueDateStr + "|\n"

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
