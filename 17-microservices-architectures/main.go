package main

import (
	"context"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      "reviews",
		Description: "KV Bucket to hold all data for the reviews microservice",
		History:     5,
	})
	if err != nil {
		log.Fatal(err)
	}

	svc, err := micro.AddService(nc, micro.Config{
		Name:        "reviews",
		Version:     "0.0.1",
		Description: "Microservices to manage ratings and reviews for various products",
		Metadata:    map[string]string{},
	})
	if err != nil {
		log.Fatal(err)
	}

	api := NewAPI(ctx, kv)
	products := svc.AddGroup("reviews.products")
	products.AddEndpoint("list_products", micro.HandlerFunc(api.ListProducts), micro.WithEndpointSubject("list"))
	products.AddEndpoint("create_product", micro.HandlerFunc(api.CreateProduct), micro.WithEndpointSubject("create"))
	products.AddEndpoint("delete_product", micro.HandlerFunc(api.CreateProduct), micro.WithEndpointSubject("delete.*"))

	runtime.Goexit()
}
