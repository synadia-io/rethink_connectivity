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
	if len(os.Args) == 1 {
		fmt.Println("Commands:")
		fmt.Println("  pub - publish to a JetStream")
		fmt.Println("  sub - subscribe to a JetStream consumer")
		os.Exit(1)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "orders",
		Description: "Stream of orders",
		Subjects:    []string{"orders.*"},
	})
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "pub":
		publish(js)
	case "sub":
		subscribe(stream)
	default:
		fmt.Println("Unknown command")
		os.Exit(1)
	}
}

func publish(js jetstream.JetStream) {
	fmt.Println("Publishing to JetStream")
	i := 0
	for {
		log.Println("Publishing message", i)
		_, err := js.Publish(context.Background(), fmt.Sprintf("orders.%d", i), []byte("Hello World"))
		if err != nil {
			log.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
		i++
	}
}

func subscribe(stream jetstream.Stream) {
	fmt.Println("Subscribing to JetStream")
	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Name:    "orders_consumer",
		Durable: "orders_consumer",
	})
	if err != nil {
		log.Fatal(err)
	}
	c, err := consumer.Consume(func(msg jetstream.Msg) {
		log.Printf("Received message: %s", msg.Subject())
		msg.Ack()
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
