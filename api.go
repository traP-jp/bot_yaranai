package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

// https://git.trap.jp/pikachu/traQ-BOT-pika-test/src/branch/main/commands/stamps.go
func getUserMessages(bot *traqwsbot.Bot, userID string, progressMessageID string) ([]traq.Message, error) {
	var messages []traq.Message
	var before = time.Now()
	res, r, err := bot.API().UserApi.GetUserStats(context.Background(), userID).Execute()
	totalmessage := res.TotalMessageCount
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	for {
		t1 := time.Now()

		res, r, err := bot.API().MessageApi.SearchMessages(context.Background()).From(userID).Limit(int32(100)).Offset(int32(0)).Before(before).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}

		fmt.Println(time.Since(t1))
		if err != nil {
			return nil, err
		}
		if len(res.Hits) == 0 {
			break
		}
		messages = append(messages, res.Hits...)
		time.Sleep(time.Millisecond * 100)
		before = messages[len(messages)-1].CreatedAt
		fmt.Println(len(messages))

		fmt.Println(len(messages))
		simpleEdit(bot, progressMessageID, strings.Repeat(":", len(messages)/100)+strings.Repeat(".", max(0, int(totalmessage)-len(messages))/100))
	}

	return messages, nil
}
func getChannelMessages(bot *traqwsbot.Bot, channelID string, progressMessageID string) ([]traq.Message, error) {
	res, r, err := bot.API().ChannelApi.GetChannelStats(context.Background(), channelID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	totalmessage := res.TotalMessageCount
	var messages []traq.Message
	var before = time.Now()
	for {
		t1 := time.Now()

		res, r, err := bot.API().MessageApi.SearchMessages(context.Background()).In(channelID).Limit(int32(100)).Offset(int32(0)).Before(before).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}

		fmt.Println(time.Since(t1))
		if err != nil {
			return nil, err
		}
		if len(res.Hits) == 0 {
			break
		}

		messages = append(messages, res.Hits...)
		time.Sleep(time.Millisecond * 100)
		before = messages[len(messages)-1].CreatedAt
		fmt.Println(len(messages))
		simpleEdit(bot, progressMessageID, strings.Repeat(":", len(messages)/100)+strings.Repeat(".", max(0, int(totalmessage)-len(messages))/100))
	}

	return messages, nil
}
func simplePost(bot *traqwsbot.Bot, c string, s string) (x string) {
	q, r, err := bot.API().
		MessageApi.
		PostMessage(context.Background(), c).
		PostMessageRequest(traq.PostMessageRequest{
			Content: s,
		}).
		Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	return q.Id
}
func simpleEdit(bot *traqwsbot.Bot, m string, s string) {
	bot.API().
		MessageApi.EditMessage(context.Background(), m).PostMessageRequest(traq.PostMessageRequest{
		Content: s,
	}).Execute()
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
