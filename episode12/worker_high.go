package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, err := nats.Connect("connect.ngs.global", nats.UserCredentials("user.creds"), nats.Name("worker_high"))
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	consumer, err := js.CreateOrUpdateConsumer(ctx, "jobs", jetstream.ConsumerConfig{
		Name:        "worker_high",
		Durable:     "worker_high",
		Description: "High priority worker pool",
		BackOff: []time.Duration{
			1 * time.Second,
			5 * time.Second,
			10 * time.Second,
		},
		MaxDeliver: 4,
	})
	if err != nil {
		log.Fatal(err)
	}
	consumer.Messages()

	c, err := consumer.Consume(func(msg jetstream.Msg) {
		meta, err := msg.Metadata()
		if err != nil {
			log.Printf("Error getting metadata: %s\n", err)
			msg.Nak()
			return
		}

		log.Printf("Not acking message message: %s\n", fmt.Sprint(meta.Sequence.Stream))
	})
	if err != nil {
		log.Fatal(err)
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	c.Stop()

}
