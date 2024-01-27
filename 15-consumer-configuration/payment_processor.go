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
		Name:        "payment_processor",
		Durable:     "payment_processor",
		Description: "Processes pending order and confirms them.",
		FilterSubjects: []string{
			"orders.*.*.pending",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	cctx, err := consumer.Consume(func(msg jetstream.Msg) {
		fmt.Println("Received message:", string(msg.Data()))
		msg.Ack()
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
