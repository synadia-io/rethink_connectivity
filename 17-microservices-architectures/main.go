package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"
)

var logger *slog.Logger
var logLevel slog.Level = slog.LevelError
var natsHandler *natsSlogHandler

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	svc, err := micro.AddService(nc, micro.Config{
		Name:        "math",
		Version:     "0.0.1",
		Description: "A simple calculator service",
	})
	if err != nil {
		log.Fatal(err)
	}

	h := slog.NewJSONHandler(&natsLogWriter{"math.log." + svc.Info().ID, nc}, nil)
	natsHandler = &natsSlogHandler{h, slog.LevelInfo}
	logger = slog.New(natsHandler)

	m := svc.AddGroup("math")

	m.AddEndpoint("add",
		micro.HandlerFunc(addHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Adds two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	m.AddEndpoint("subtract",
		micro.HandlerFunc(subtractHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Subtracts two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	m.AddEndpoint("multiply",
		micro.HandlerFunc(multiplyHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Multiplies two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	m.AddEndpoint("divide",
		micro.HandlerFunc(divideHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Divides two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	log.Printf("Listening on %q for %q, %q, %q, %q \n", nc.ConnectedAddr(), "math.add", "math.subtract", "math.multiply", "math.divide")

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      "config",
		Description: "config for my microservices",
	})
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := kv.Watch(ctx, "math.log_level")
	if err != nil {
		log.Fatal(err)
	}

	for entry := range watcher.Updates() {
		if entry == nil {
			continue
		}

		level, err := strconv.Atoi(string(entry.Value()))
		if err != nil {
			logger.Error("invalid log level", "error", err)
			continue
		}

		logLevel = slog.Level(level)
		natsHandler.level = logLevel
		logger.Log(ctx, 6, "log level set to "+logLevel.String())
	}

}

func addHandler(req micro.Request) {
	logger.Info("addHandler called", "subject", req.Subject())
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		logger.Error("error", "error", err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A + calcRequest.B,
	})
}

func subtractHandler(req micro.Request) {
	logger.Info("subtractHandler called", "subject", req.Subject())
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		logger.Error("error", "error", err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A - calcRequest.B,
	})
}

func multiplyHandler(req micro.Request) {
	logger.Info("multiplyHandler called", "subject", req.Subject())
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		logger.Error("error", "error", err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A * calcRequest.B,
	})
}

func divideHandler(req micro.Request) {
	logger.Info("divideHandler called", "subject", req.Subject())
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		logger.Error("error", "error", err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	if calcRequest.B == 0 {
		logger.Error("cannot divide by zero", "subject", req.Subject())
		req.Error("400", "cannot divide by zero", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A / calcRequest.B,
	})
}
