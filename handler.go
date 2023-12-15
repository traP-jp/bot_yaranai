package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func checkHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	resp, err := getChannelMessages(bot, p.Message.ChannelID, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	PostMessagesWithStamp(bot, resp, c, 1)
}
func checkUserHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	userlist, _, _ := bot.API().UserApi.GetUsers(context.Background()).Execute()
	cmd := strings.Split(p.Message.Text, " ")
	userID := p.Message.User.ID
	if len(cmd) >= 5 {
		for _, v := range userlist {
			if v.Name == cmd[4] {
				userID = v.Id
			}
		}
	}
	resp, err := getUserMessages(bot, userID, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	atleast := 1
	if len(cmd) >= 6 {
		atleast, _ = strconv.Atoi(cmd[5])
	}
	PostMessagesWithStamp(bot, resp, c, atleast)
}

func heatMapHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	resp, err := getUserMessages(bot, p.Message.User.ID, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	daily := map[string]int{}
	for _, v := range resp {
		daily[v.CreatedAt.Format("2006-01-02")] += 1
	}
	s := ""
	for i := 0; i < 7; i++ {
		s += fmt.Sprintf("%s : %d\n", time.Now().AddDate(0, 0, -i).Format("2006-01-02"), daily[time.Now().AddDate(0, 0, -i).Format("2006-01-02")])
	}
	if len(s) > 3000 {
		s = s[:3000] + "\n(snip)"
	}
	simpleEdit(bot, c, s)
}
func stampCountHandrer(bot *traqwsbot.Bot, p *payload.MessageCreated) {
	c := simplePost(bot, p.Message.ChannelID, "実行中...")
	userlist, _, _ := bot.API().UserApi.GetUsers(context.Background()).Execute()
	cmd := strings.Split(p.Message.Text, " ")
	userID := p.Message.User.ID
	if len(cmd) >= 4 {
		for _, v := range userlist {
			if v.Name == cmd[3] {
				userID = v.Id
			}
		}
	}
	resp, err := getUserMessages(bot, userID, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	//-------------------------------------
	target := "w"
	if len(cmd) >= 5 {
		target = cmd[4]
	}
	stamplist, r, err := bot.API().StampApi.GetStamps(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", r)
	}
	for _, stamp := range stamplist {
		if target == stamp.Name || target == ":"+stamp.Name+":" {
			target = stamp.Id
		}
	}
	fmt.Println(target)
	count := 0
	total := 0
	for _, v := range resp {
		for _, w := range v.Stamps {
			if w.StampId == target {
				count += 1
				total += int(w.Count)
			}
		}
	}
	userstat, r, err := bot.API().UserApi.GetUserStats(context.Background(), userID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", r)
	}
	sendcount, sendtotal := 0, 0
	for _, v := range userstat.Stamps {
		if v.Id == target {
			sendcount = int(v.Count)
			sendtotal = int(v.Total)
		}
	}
	s := "received: " + strconv.Itoa(count) + "(total: " + strconv.Itoa(total) + ")\n"
	s += "sent: " + strconv.Itoa(sendcount) + "(total: " + strconv.Itoa(sendtotal) + ")\n"
	simpleEdit(bot, c, s)
}
func PostMessagesWithStamp(bot *traqwsbot.Bot, resp []traq.Message, c string, atleast int) {
	stamplist, r, err := bot.API().StampApi.GetStamps(context.Background()).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	s := ""
	for _, v := range resp {
		c := ""
		f := 0
		q := []rune(v.Content)
		if len(q) > 50 {
			c += string(q[0:50]) + "... : "
		} else {
			c += string(q) + " : "
		}
		for _, w := range v.Stamps {
			for _, stamp := range stamplist {
				if w.StampId == stamp.Id {
					var i int32
					for i = 0; i < w.Count; i++ {
						c += ":" + stamp.Name + ":"
					}
					f += 1
				}
			}
		}
		c += "\n"
		if f >= atleast {
			s += c
		}
	}
	if len(s) > 3000 {
		s = s[:3000] + "\n(snip)"
	}
	simpleEdit(bot, c, s)
}
