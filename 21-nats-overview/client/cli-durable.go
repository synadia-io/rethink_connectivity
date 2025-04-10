package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

func RunCLIDurable() {
	// Connect to NATS
  fmt.Println(nats.DefaultURL)
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer nc.Close()

  js, err := nc.JetStream()
  if err != nil {
    fmt.Println(err)
  }

  _, err = js.AddStream(&nats.StreamConfig{
    Name: "CHAT",
    Subjects: []string{"chat.>"},
  })
  if err != nil {
    fmt.Println(err)
  }

	// Subscribe to chat messages
	nc.Subscribe("chat.messages", func(m *nats.Msg) {
		fmt.Printf("Received: %s\n", string(m.Data))
	})

  // Retrieve chat history
  consumer, err := js.PullSubscribe("chat.messages", "")
  if err != nil {
    fmt.Println(err)
    return
  }

  msgs, err := consumer.Fetch(100)
  if err != nil {
    fmt.Println(err)
  } else {
    for _, msg := range msgs {
      fmt.Printf("Received: %s\n", string(msg.Data))
    }
  }

	// Publish chat messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.ToLower(msg) == "quit" {
			break
		}
		js.Publish("chat.messages", []byte(msg))
	}
}
