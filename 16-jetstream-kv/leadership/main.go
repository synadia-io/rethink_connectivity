package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go/jetstream"
)

var isLeader = false
var name string
var rev uint64

func main() {
	name = os.Args[1]
	ctx := context.Background()
	log.SetPrefix(fmt.Sprintf("[%s] ", name))
	nc, err := natscontext.Connect(natscontext.SelectedContext())
	if err != nil {
		log.Fatalln(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalln(err)
	}

	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:   "leadership",
		TTL:      time.Second * 3,
		MaxBytes: 1024 * 1024,
	})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		time.Sleep(1 * time.Second)

		if isLeader {
			rev, err = kv.Update(ctx, "leader", []byte(name), rev)
			if err != nil {
				log.Println("Lost leadership:", err)
				isLeader = false
			}
			continue
		} else {
			log.Println("Trying to assume leadership...")
			rev, err = kv.Create(ctx, "leader", []byte(name))
			if err != nil {
				continue
			}

			log.Println("Became Leader")
			isLeader = true
		}
	}
}
