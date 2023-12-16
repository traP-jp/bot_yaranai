package main

import (
	"fmt"
	"log"
	"os"
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
		if cmd[1] == "task" {
			getTest(bot, p.Message.ChannelID)
		} else if cmd[1] == "help" {
			bytes, err := os.ReadFile("help.txt")
			if err != nil {
				panic(err)
			}
			simplePost(bot, p.Message.ChannelID, string(bytes))
		} else if cmd[1] == "condition" {
			if len(cmd) == 2 {
				bytes, err := os.ReadFile("help.txt")
				if err != nil {
					panic(err)
				}
				simplePost(bot, p.Message.ChannelID, string(bytes))
			} else {
				switch cmd[2] {
				case "list":
					debuggetCondition(bot, p.Message.ChannelID)
				case "add":
					if len(cmd) == 3 {
						simplePost(bot, p.Message.ChannelID, "Name cannot be empty")
					} else {
						postCondition(bot, p.Message.ChannelID, strings.Join(cmd[3:]," "))
					}
				case "delete":
					if len(cmd) != 4 {
						simplePost(bot, p.Message.ChannelID, "Please specify a condition_id")
					} else {
						deleteCondition(bot, p.Message.ChannelID, cmd[3])
					}

				default:
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
