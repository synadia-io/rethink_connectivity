package main

import (
	"fmt"
	"os"

	"nats-chat/client"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [cli|web]")
		return
	}

	clientType := os.Args[1]
	switch clientType {
	case "cli":
		client.RunCLI()
	case "cli-durable":
		client.RunCLIDurable()
	case "web":
		client.RunWeb()
	case "web-durable":
		client.RunWebDurable()
	case "twitch":
		if len(os.Args) < 3 {
			fmt.Println("Please provide channel name for twitch")
			return
		}
		channelName := os.Args[2]
		client.RunTwitch(channelName)
		fmt.Printf("Connecting to channel: %s\n", channelName)
	default:
		fmt.Printf("Unknown client type: %s\n", clientType)
		fmt.Println("Usage: go run main.go [cli|web]")
	}
}
