package main

import (
	"context"
	"log"

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

	js.CreateOrUpdateConsumer(ctx, "orders", jetstream.ConsumerConfig{
		Name:               "",
		Durable:            "",
		Description:        "",
		DeliverPolicy:      0,
		OptStartSeq:        0,
		OptStartTime:       &time.Time{},
		AckPolicy:          0,
		AckWait:            0,
		MaxDeliver:         0,
		BackOff:            []time.Duration{},
		FilterSubject:      "",
		ReplayPolicy:       0,
		RateLimit:          0,
		SampleFrequency:    "",
		MaxWaiting:         0,
		MaxAckPending:      0,
		HeadersOnly:        false,
		MaxRequestBatch:    0,
		MaxRequestExpires:  0,
		MaxRequestMaxBytes: 0,
		InactiveThreshold:  0,
		Replicas:           0,
		MemoryStorage:      false,
		FilterSubjects:     []string{},
		Metadata:           map[string]string{},
	})
}
