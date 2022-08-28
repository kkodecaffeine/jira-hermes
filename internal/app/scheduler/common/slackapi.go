package common

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

func OpenConversation(email string) string {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		panic("NO SLACK_TOKEN")
	}

	client := slack.New(token)
	userInfo, err := client.GetUserByEmail(email)
	if err != nil {
		fmt.Printf("Error user not found: %s\n", email)
		return ""
	}

	channel, _, _, err := client.OpenConversation(&slack.OpenConversationParameters{Users: []string{userInfo.ID}})
	if err != nil {
		fmt.Printf("Error opening IM: %s\n", err)
	}

	return channel.ID
}

func PostMessage(message string, target string) (string, string) {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		panic("NO SLACK_TOKEN")
	}

	channelID := os.Getenv("SLACK_CHANNEL_ID")
	if channelID == "" {
		panic("NO SLACK_CHANNEL_ID")
	}

	if target != "" {
		channelID = target
	}

	client := slack.New(token)
	attachment := slack.Attachment{
		Pretext: message,
	}

	channel, timestamp, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		fmt.Println("post " + err.Error())
	}

	return channel, timestamp
}

func UpdateMessage(channelID string, ts string, message string) string {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		panic("NO SLACK_TOKEN")
	}

	client := slack.New(token)
	attachment := slack.Attachment{
		Text: message,
	}

	_, _, timestamp, err := client.UpdateMessage(
		channelID,
		ts,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		fmt.Println("update " + err.Error())
	}

	return timestamp
}

func ReplyMessage(channelID string, ts string, message slack.Message) {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		panic("NO SLACK_TOKEN")
	}

	client := slack.New(token)
	attachment := slack.Attachment{
		Blocks: message.Blocks,
	}

	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionTS(ts),
	)

	if err != nil {
		fmt.Println("reply " + err.Error())
	}
}
