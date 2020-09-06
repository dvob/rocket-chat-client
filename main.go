package main

import (
	"fmt"
	"os"

	"github.com/dsbrng25b/rocket-chat-client/pkg/rocketchat"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing action")
		os.Exit(1)
	}

	rcURL := os.Getenv("RC_URL")
	rcUserID := os.Getenv("RC_USER_ID")
	rcUserToken := os.Getenv("RC_USER_TOKEN")

	client := rocketchat.NewClient(rcURL, rcUserID, rcUserToken)

	action := os.Args[1]
	switch action {
	case "list":
		channels, err := client.ListChannels()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, channel := range channels {
			fmt.Println("#" + channel.Name)
		}
		users, err := client.ListUsers()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, user := range users {
			fmt.Println("@" + user.Username)
		}

	case "send":
		if len(os.Args) < 4 {
			fmt.Println("missing arguments")
			os.Exit(1)
		}
		dest := os.Args[2]
		text := os.Args[3]

		client.SendMessage(dest, text, "")
	default:
		fmt.Println("unknown action:", action)
		os.Exit(1)
	}
}
