package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type JSONSource struct {
	BotToken string    `json:"bot_token"`
	Messages []Message `json:"channels"`
}

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type ChannelIDByName map[string]string

func GetChannelIDs(client *slack.Client) (ChannelIDByName, error) {
	conversations, _, err := client.GetConversations(&slack.GetConversationsParameters{TeamID: "T04A9BS7UJF"})
	if err != nil {
		return nil, err
	}

	channelIDByName := ChannelIDByName{}

	for _, conversation := range conversations {
		id := conversation.ID
		name := conversation.Name
		channelIDByName[name] = id
	}

	return channelIDByName, nil
}

var filename string

func init() {
	flag.StringVar(&filename, "f", "messages.json", "define a json file to pass")

	flag.Parse()
}

func unmarshalMessagesFileJSON(filename string) (*JSONSource, error) {
	var jsonHandler JSONSource

	temp, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(temp, &jsonHandler)
	if err != nil {
		return nil, err
	}

	return &jsonHandler, nil
}

func main() {
	jsonData, err := unmarshalMessagesFileJSON(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := slack.New(jsonData.BotToken)

	channelIDByName, err := GetChannelIDs(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, message := range jsonData.Messages {
		_, timestamp, err := client.PostMessage(
			channelIDByName[message.Channel],
			slack.MsgOptionText(message.Text, false),
		)

		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Sent to channel %s at %s\n", message.Channel, timestamp)
	}
}