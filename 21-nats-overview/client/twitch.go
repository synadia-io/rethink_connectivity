package client

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v4"
  "github.com/nats-io/nats.go"
)

func RunTwitch(channel string) {
  nc, err := nats.Connect(nats.DefaultURL)
  if err != nil {
    fmt.Println(err)
  }
  defer nc.Close()

  js, err := nc.JetStream()
  if err != nil {
    fmt.Println(err)
  }

	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
    chatMessage := fmt.Sprintf("%s: %s (twitch)", message.User.DisplayName, message.Message)

    _, err := js.Publish("chat.messages", []byte(chatMessage))
    if err != nil {
      fmt.Println("Error publishing to NATS:", err)
    }
	})

	client.Join(channel)

  err = client.Connect()
  if err != nil {
    panic(err)
  }
}
