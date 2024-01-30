package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// cctx, err := eventSourcingConsumer(js)
	// cctx, err := reportingConsumer(js)
	cctx, err := lookupConsumer(js)
	if err != nil {
		log.Fatal(err)
	}
	defer cctx.Stop()

	// Gracefully shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}

func eventSourcingConsumer(js jetstream.JetStream) (jetstream.ConsumeContext, error) {
	ctx := context.Background()

	consumer, err := js.CreateOrUpdateConsumer(ctx, "orders", jetstream.ConsumerConfig{
		Name:          "event-sourcing",
		Durable:       "event-sourcing",
		MaxAckPending: 1,
		Description:   "Processes pending orders in strict ordering (slow!!!)",
		FilterSubjects: []string{
			"orders.>",
		},
	})
	if err != nil {
		return nil, err
	}

	return consumer.Consume(func(msg jetstream.Msg) {
		meta, err := msg.Metadata()
		if err != nil {
			fmt.Println("Error getting metadata:", err)
			return
		}
		fmt.Println("Processing message:", meta.Sequence.Consumer, msg.Subject())
		msg.Ack()
	})
}

func reportingConsumer(js jetstream.JetStream) (jetstream.ConsumeContext, error) {
	ctx := context.Background()

	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)

	consumer, err := js.CreateOrUpdateConsumer(ctx, "orders", jetstream.ConsumerConfig{
		Name:            "reporting",
		HeadersOnly:     true,
		MaxRequestBatch: 1000,
		AckPolicy:       jetstream.AckNonePolicy,
		DeliverPolicy:   jetstream.DeliverByStartTimePolicy,
		OptStartTime:    &sevenDaysAgo,
		FilterSubjects: []string{
			"orders.US.*.shipped",
		},
	})
	if err != nil {
		return nil, err
	}

	count := 0
	return consumer.Consume(func(msg jetstream.Msg) {
		count++
		fmt.Println("Count", count, msg.Subject())
	}, jetstream.PullMaxMessages(1000))
}

func lookupConsumer(js jetstream.JetStream) (jetstream.ConsumeContext, error) {
	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, "orders", jetstream.ConsumerConfig{
		Name:          "lookup",
		AckPolicy:     jetstream.AckNonePolicy,
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		FilterSubjects: []string{
			"orders.US.6080417374930863726.*",
		},
	})
	if err != nil {
		return nil, err
	}

	return consumer.Consume(func(msg jetstream.Msg) {
		fmt.Println(msg.Subject(), string(msg.Data()))
	})
}
