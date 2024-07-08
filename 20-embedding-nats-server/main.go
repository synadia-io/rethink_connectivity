package main

import (
	"log"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"time"
)

func main() {

	opts := &server.Options{}

	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	ns.ConfigureLogger()
	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		log.Fatal("Embedded nats-server timed out")
	}

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe("hello.world", func(msg *nats.Msg) {
		log.Println("Hello world")
	})

	ns.WaitForShutdown()
}
