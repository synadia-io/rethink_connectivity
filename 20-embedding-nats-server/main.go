package main

import (
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"time"
)

func main() {
}

func RunEmbeddedServer(inProcess bool, enableLogging bool) (*nats.Conn, *server.Server, error) {
	opts := &server.Options{
		DontListen: inProcess,
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, nil, err
	}

	if enableLogging {
		ns.ConfigureLogger()
	}
	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, nil, err
	}

	clientOpts := []nats.Option{}
	if inProcess {
		clientOpts = append(clientOpts, nats.InProcessServer(ns))
	}

	nc, err := nats.Connect(nats.DefaultURL, clientOpts...)
	if err != nil {
		return nil, nil, err
	}

	return nc, ns, err
}
