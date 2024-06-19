package main

import (
	"fmt"
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	RunAuthService()
}

func RunAuthService() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}

	_, err = micro.AddService(nc, micro.Config{
		Name:        "auth",
		Version:     "0.0.1",
		Description: "Handle authorization of Google JWTs for chat applications",
		Endpoint: &micro.EndpointConfig{
			Subject: "auth",
			Handler: micro.HandlerFunc(func(r micro.Request) {
				fmt.Println("Received Request", r.Data())
				r.Error("500", "Not implemented", nil)
			}),
		},
	})
	if err != nil {
		return err
	}

	runtime.Goexit()
	return nil
}
