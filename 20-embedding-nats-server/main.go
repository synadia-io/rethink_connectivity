package main

import (
	"log"
	"net/url"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, ns, err := RunEmbeddedServer(true, true)
	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe("hello.world", func(msg *nats.Msg) {
		log.Println("message received!")
		msg.Respond([]byte("Ahoy there!"))
	})

	ns.WaitForShutdown()
}

func RunEmbeddedServer(inProcess bool, enableLogging bool) (*nats.Conn, *server.Server, error) {
	leafURL, err := url.Parse("nats-leaf://connect.ngs.global")
	// leafURL, err := url.Parse("nats-leaf://0.0.0.0:7422")
	if err != nil {
		return nil, nil, err
	}

	opts := &server.Options{
		ServerName:      "embedded_server",
		DontListen:      inProcess,
		JetStream:       true,
		JetStreamDomain: "embedded",
		LeafNode: server.LeafNodeOpts{
			Remotes: []*server.RemoteLeafOpts{
				{
					URLs:        []*url.URL{leafURL},
					Credentials: "./leafnode.creds",
				},
			},
		},
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
