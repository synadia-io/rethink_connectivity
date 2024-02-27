package main

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/redis/go-redis/v9"
)

var (
	js jetstream.JetStream
	kv jetstream.KeyValue
	nc *nats.Conn

	rdb *redis.Client
)

func setup() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var err error
	nc, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}

	js, err = jetstream.New(nc)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	kv, err = js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      "benchmark",
		Description: "A stream used for benchmarking",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func shutdown() {
	js.DeleteKeyValue(context.Background(), "benchmark")
	nc.Drain()
	rdb.FlushAll(context.Background())
	rdb.Close()

}

func put(key, value string) error {
	_, err := kv.Put(context.Background(), key, []byte(value))
	return err
}

func set(key, value string) error {
	return rdb.Set(context.Background(), key, value, 0).Err()
}
