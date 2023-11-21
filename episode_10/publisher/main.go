package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, err := nats.Connect("connect.ngs.global", nats.UserCredentials("../user.creds"), nats.Name("Orders Publisher"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ctx := context.Background()
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	// Create a stream
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "orders",
		Description: "Messages for orders",
		Subjects: []string{
			"orders.>",
		},
		MaxBytes: 1024 * 1024 * 1024,
	})
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for {
		i++
		_, err = js.Publish(ctx, fmt.Sprintf("orders.%d", i), []byte("Hello World"))
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Published message %d", i)
	}
}
