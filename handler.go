package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func getTask(bot *traqwsbot.Bot, userID string, channelID string) {
	var tasks []Task
	if err := db.Select(&tasks, "SELECT * FROM `task` WHERE user = ?", userID); err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "There is no such task of yours")
		return
	}

	res := "## タスク一覧\n|ID|タスク名|詳細|状況|ウェイト|期限|\n|---|---|---|---|---|---|\n"
	for _, v := range tasks {
		idStr := strconv.Itoa(v.Id)

		//与えられたcondition_idからcondition名を取得する
		var conditionName string
		err := db.Get(&conditionName, "SELECT `condition` FROM `condition` WHERE `condition_id`=?", v.ConditionId)
		if err != nil {
			fmt.Println(err)
			simplePost(bot, channelID, "Condition get error")
			return
		}
		weightstr := strconv.Itoa(v.Difficulty)
		dueDateStr := v.DueDate.Format("2006-01-02")
		res += "|" + idStr + "|" + v.Title + "|" + v.Description + "|" + conditionName + "|" + weightstr + "|" + dueDateStr + "|\n"
	}
	simplePost(bot, channelID, res)
}

func postTask(bot *traqwsbot.Bot, userID string, channelID string, newTask TaskWithoutId, conditionName string) {

	//与えられたcondition名からcondition_idを取得する
	var conditionId int
	err := db.Get(&conditionId, "SELECT `condition_id` FROM `condition` WHERE `condition`=?", conditionName)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		simplePost(bot, channelID, "There is no such condition")
		return
	} else if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Condition get error")
		return
	}

	var dateOfNow = time.Now().Format("2006-01-02")
	res, err := db.Exec("INSERT INTO task (user, title, description, condition_id, difficulty, created_at, updated_at, dueDate) VALUES (?,?,?,?,?,?,?,?)", userID, newTask.Title, newTask.Description, conditionId, newTask.Difficulty, dateOfNow, dateOfNow, newTask.DueDate)

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

	//与えられたcondition_idからcondition名を取得する

	weightstr := strconv.Itoa(newTask.Difficulty)

	resStr := "## タスクを追加しました！\n|ID|タスク名|詳細|状況|ウェイト|期限|\n|---|---|---|---|---|---|\n|" + strconv.Itoa(int(taskID)) + "|" + newTask.Title + "|" + newTask.Description + "|" + conditionName + "|" + weightstr + "|" + newTask.DueDate + "|\n"
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

	information := [...]string{"title", "description", "conditionName", "difficulty", "dueDate"}
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
			} else if information[i] == "conditionName" {
				//与えられたcondition名からcondition_idを取得する
				var conditionId int
				err := db.Get(&conditionId, "SELECT `condition_id` FROM `condition` WHERE `condition`=?", change)
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println(err)
					simplePost(bot, channelID, "There is no such condition")
					return
				} else if err != nil {
					fmt.Println(err)
					simplePost(bot, channelID, "Condition get error")
					return
				}
				_, err = db.Exec("UPDATE task SET condition_id = ? WHERE id = ? AND user = ?", conditionId, taskID, userID)
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

	//与えられたcondition_idからcondition名を取得する
	var conditionName string
	err = db.Get(&conditionName, "SELECT `condition` FROM `condition` WHERE `condition_id`=?", updatedTask.ConditionId)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Condition get error")
		return
	}

	difficultyInt := strconv.Itoa(updatedTask.Difficulty)

	dueDateStr := updatedTask.DueDate.Format("2006-01-02")

	resStr := "## タスクを変更しました。\n変更結果\n|タスク名|詳細|状況|ウェイト|期限|\n|---|---|---|---|---|\n|" + updatedTask.Title + "|" + updatedTask.Description + "|" + conditionName + "|" + difficultyInt + "|" + dueDateStr + "|\n"

	simplePost(bot, channelID, resStr)
}

func deleteTask(bot *traqwsbot.Bot, userID string, channelID string, taskIDstr string) {
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

	//与えられたcondition_idからcondition名を取得する
	var conditionName string
	err = db.Get(&conditionName, "SELECT `condition` FROM `condition` WHERE `condition_id`=?", task.ConditionId)
	if err != nil {
		fmt.Println(err)
		simplePost(bot, channelID, "Condition get error")
		return
	}

	weightstr := strconv.Itoa(task.Difficulty)
	dueDateStr := task.DueDate.Format("2006-01-02")

	//ログの出力
	simplePost(bot, channelID, "タスクの消去が完了しました。消去内容は以下の通りです\n| ID | タスク名 |詳細|状況|ウェイト|期限|\n| --- | --- | --- | --- | --- | --- |\n| "+taskIDstr+" | "+task.Title+" |"+task.Description+" |"+conditionName+" |"+weightstr+" |"+dueDateStr+" |\n ")
}
