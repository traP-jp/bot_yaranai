package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

var (
	db *sqlx.DB
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("MARIADB_USER") == "" {
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(os.Getenv("MARIADB_USER"))
	fmt.Println("aa")
	conf := mysql.Config{
		User:                 os.Getenv("MARIADB_USER"),
		Passwd:               os.Getenv("MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MARIADB_HOSTNAME") + ":" + os.Getenv("MARIADB_PORT"),
		DBName:               os.Getenv("MARIADB_DATABASE"),
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	_db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("conntected")
	db = _db
	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		fmt.Println(p.Message.Text)
		cmd := strings.Split(p.Message.Text, " ")

		userID := p.Message.User.ID

		if cmd[1] == "task" {
			if len(cmd) == 2 {
				getTest(bot, p.Message.ChannelID)
			} else {
				switch cmd[2] {
				case "get":
					getTask(bot, userID, p.Message.ChannelID)
				case "post":
					conditionIdInt, err := strconv.Atoi(cmd[5])
					if err != nil {
						fmt.Println(err)
						simplePost(bot, p.Message.ChannelID, "Condition ID は数値にしてください")
						return
					}
					difficultyInt, err := strconv.Atoi(cmd[6])
					if err != nil {
						fmt.Println(err)
						simplePost(bot, p.Message.ChannelID, "Difficulty は数値にしてください")
						return
					}
					newTask := TaskWithoutId{
						Title:       cmd[3],
						Description: cmd[4],
						ConditionId: conditionIdInt,
						Difficulty:  difficultyInt,
						DueDate:     cmd[7],
					}
					postTask(bot, userID, p.Message.ChannelID, newTask)
				case "edit": //タスクの編集(PUT: /task に相当)
					taskId, err := strconv.Atoi(cmd[3])
					if err != nil {
						fmt.Println(err)
						simplePost(bot, p.Message.ChannelID, "タスク ID は数値を入力してください。")
						return
					}
					if len(cmd) != 9 {
						if len(cmd) > 9 {
							simplePost(bot, p.Message.ChannelID, "入力が多すぎます")
						} else {
							simplePost(bot, p.Message.ChannelID, "入力が少なすぎます")
						}
						return
					}
					var changeList = cmd[4:]
					putTask(bot, taskId, userID, p.Message.ChannelID, [5]string(changeList))
				case "delete": //コンディションの消去(DELETE: /condition に相当)
					if len(cmd) != 4 {
						simplePost(bot, p.Message.ChannelID, "Please specify a taskid")
					} else {
						deleteTask(bot, userID, p.Message.ChannelID, cmd[3])
					}
				default:
					simplePost(bot, p.Message.ChannelID, "No such command")
				}
			}
		} else if cmd[1] == "help" {
			bytes, err := os.ReadFile("help.txt")
			if err != nil {
				panic(err)
			}
			simplePost(bot, p.Message.ChannelID, string(bytes))
		} else if cmd[1] == "condition" {
			//単にconditionと打たれただけならヘルプを表示
			if len(cmd) == 2 {
				bytes, err := os.ReadFile("help.txt")
				if err != nil {
					panic(err)
				}
				simplePost(bot, p.Message.ChannelID, string(bytes))
			} else {
				switch cmd[2] {
				case "get": //ユーザー毎コンディションリストの取得
					getCondition(bot, p.Message.ChannelID, userID)
				case "post": //コンディションの追加(POST: /condition に相当)
					//引数不足の場合
					if len(cmd) == 3 {
						simplePost(bot, p.Message.ChannelID, "Name cannot be empty")
					} else {
						postCondition(bot, p.Message.ChannelID, strings.Join(cmd[3:], " "), userID)
					}
				case "edit": //コンディションの編集(PUT: /condition に相当)
					//引数不足の場合(id不足, name不足 入替パターンについてはidが数値変換できなかった場合のエラー(ハンドラ内)で拾う)
					if len(cmd) == 3 {
						simplePost(bot, p.Message.ChannelID, "Please specify a condition_id")
					} else if len(cmd) == 4 {
						simplePost(bot, p.Message.ChannelID, "Name cannot be empty")
					} else {
						putCondition(bot, p.Message.ChannelID, cmd[3], strings.Join(cmd[4:], " "), userID)
					}
				case "delete": //コンディションの消去(DELETE: /condition に相当)
					if len(cmd) != 4 {
						simplePost(bot, p.Message.ChannelID, "Please specify a condition_id")
					} else {
						deleteCondition(bot, p.Message.ChannelID, cmd[3], userID)
					}

				default: //存在しないコマンドの場合
					bytes, err := os.ReadFile("help.txt")
					if err != nil {
						panic(err)
					}
					simplePost(bot, p.Message.ChannelID, "No such command\n"+string(bytes))
				}
			}
		} else {
			simplePost(bot, p.Message.ChannelID, "No such command")
		}

	})

	if err := bot.Start(); err != nil {
		panic(err)
	}
}
