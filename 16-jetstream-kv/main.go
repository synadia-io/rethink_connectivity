package main

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx := context.Background()
	nc, err := natscontext.Connect(natscontext.SelectedContext())
	if err != nil {
		log.Fatalln(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalln(err)
	}

	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      "my_kv_bucket",
		Description: "My cool key value store",
		MaxBytes:    1024 * 1024 * 1024, // 1GB
	})
	if err != nil {
		log.Fatalln(err)
	}

	// get a value from the kv store every 1 second
	for {
		t := time.Now()
		v, err := kv.Get(ctx, "hello.world")
		if err != nil {
			log.Println(err)
		} else {
			elapsed := time.Since(t)
			log.Printf("Got value: %s in %s\n", v.Value(), elapsed)
		}
		time.Sleep(time.Second * 1)
	}

}
