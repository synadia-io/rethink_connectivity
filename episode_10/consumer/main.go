package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, err := nats.Connect("connect.ngs.global", nats.UserCredentials("../user.creds"), nats.Name("Orders Consumer"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ctx := context.Background()
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := js.Stream(ctx, "orders")
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    "order_processor",
		Durable: "order_processor",
	})
	if err != nil {
		log.Fatal(err)
	}

	cctx, err := consumer.Consume(func(msg jetstream.Msg) {
		log.Printf("Received message: %s", string(msg.Subject()))
		msg.Ack()
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cctx.Stop()

	// gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
