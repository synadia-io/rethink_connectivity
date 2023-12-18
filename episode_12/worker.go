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
	priority, id := os.Args[1], os.Args[2]
	consumerName := fmt.Sprintf("worker_%s", priority)
	workerName := fmt.Sprintf("worker_%s_%s", priority, id)

	log.Default().SetPrefix(fmt.Sprintf("[%s] ", workerName))

	nc, err := nats.Connect("connect.ngs.global", nats.UserCredentials(consumerName+".creds"), nats.Name(workerName))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	log.Printf("connected to %s", nc.ConnectedUrl())

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	consumer, err := js.CreateOrUpdateConsumer(ctx, "jobs", jetstream.ConsumerConfig{
		Name:        consumerName,
		Durable:     consumerName,
		Description: fmt.Sprintf("Worker pool with priority %s", priority),
		BackOff: []time.Duration{
			5 * time.Second,
			10 * time.Second,
			15 * time.Second,
		},
		MaxDeliver: 4,
		FilterSubjects: []string{
			fmt.Sprintf("jobs.%s.>", priority),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	consumer.Messages()

	c, err := consumer.Consume(func(msg jetstream.Msg) {
		meta, err := msg.Metadata()
		if err != nil {
			log.Printf("Error getting metadata: %s\n", err)
			return
		}
		log.Printf("Received message sequence: %d\n", meta.Sequence.Stream)
		time.Sleep(10 * time.Millisecond)
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
