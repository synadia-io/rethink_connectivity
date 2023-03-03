package main

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/cdipaolo/sentiment"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}

func Connect() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}

	model, err := sentiment.Restore()
	if err != nil {
		return err
	}

	schema, err := Schema()
	if err != nil {
		return err
	}

	micro.AddService(nc, micro.Config{
		Name:        "sentiment",
		Description: "perform sentiment analysis on a particular piece of of text",
		Version:     "0.0.1",
		Endpoint: &micro.EndpointConfig{
			Subject: "sentiment",
			Schema:  schema,
			Handler: micro.HandlerFunc(func(req micro.Request) {
				var r SentimentRequest
				err := json.Unmarshal(req.Data(), &r)
				if err != nil {
					req.Error("400", err.Error(), nil)
					return
				}

				req.RespondJSON(AnalyzeSentiment(&r, model))
			}),
		},
	})

	log.Println("Listening on 'sentiment'", nc.ConnectedAddr())

	return nil
}
