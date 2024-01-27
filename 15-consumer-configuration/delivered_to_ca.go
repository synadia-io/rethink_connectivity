package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, err := nats.Connect("connect.ngs.global", nats.UserCredentials("user.creds"), nats.Name("My Consumer"))
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	consumer, err := js.CreateOrUpdateConsumer(ctx, "orders", jetstream.ConsumerConfig{
		Name:        "reporting",
		Description: "How many orders have been delivered to CA?",
		AckPolicy:   jetstream.AckNonePolicy,
		HeadersOnly: true,
		MaxDeliver:  100000,
		FilterSubjects: []string{
			"orders.CA.*.delivered",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	cctx, err := consumer.Consume(func(msg jetstream.Msg) {
		count++
		fmt.Println("Total orders delivered to CA:", count)
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cctx.Stop()

	// Gracefully shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}
