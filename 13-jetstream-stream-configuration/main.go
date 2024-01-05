package main

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:        "",
		Description: "",
		Storage:     0,

		Retention:            0,
		MaxConsumers:         0,
		MaxMsgs:              0,
		MaxBytes:             0,
		Discard:              0,
		DiscardNewPerSubject: false,
		MaxAge:               0,
		MaxMsgsPerSubject:    0,
		MaxMsgSize:           0,
		Replicas:             0,
		NoAck:                false,
		Template:             "",
		Duplicates:           0,
		Placement:            &jetstream.Placement{},
		Mirror:               &jetstream.StreamSource{},
		Sources:              []*jetstream.StreamSource{},
		Sealed:               false,
		DenyDelete:           false,
		DenyPurge:            false,
		AllowRollup:          false,
		Compression:          0,
		FirstSeq:             0,
		SubjectTransform:     &jetstream.SubjectTransformConfig{},
		RePublish:            &jetstream.RePublish{},
		AllowDirect:          false,
		MirrorDirect:         false,
		ConsumerLimits:       jetstream.StreamConsumerLimits{},
		Metadata:             map[string]string{},
	})

}
