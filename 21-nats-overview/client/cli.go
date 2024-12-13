package client
import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

func RunCLI() {
	// Connect to NATS
  fmt.Println(nats.DefaultURL)
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer nc.Close()

	// Subscribe to chat messages
	nc.Subscribe("chat", func(m *nats.Msg) {
		fmt.Printf("Received: %s\n", string(m.Data))
	})

	// Publish chat messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.ToLower(msg) == "quit" {
			break
		}
		nc.Publish("chat", []byte(msg))
	}
}
