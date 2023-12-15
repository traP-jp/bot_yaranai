package main

import (
	"context"
	"fmt"
	"os"

	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

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
