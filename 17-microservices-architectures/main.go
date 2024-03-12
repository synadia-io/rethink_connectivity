package main

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

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

	math := svc.AddGroup("math")

	math.AddEndpoint("add",
		micro.HandlerFunc(addHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Adds two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	math.AddEndpoint("subtract",
		micro.HandlerFunc(addHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Subtracts two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	math.AddEndpoint("multiply",
		micro.HandlerFunc(addHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Multiplies two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	math.AddEndpoint("divide",
		micro.HandlerFunc(addHandler),
		micro.WithEndpointMetadata(map[string]string{
			"description":     "Divides two numbers",
			"format":          "application/json",
			"request_schema":  SchemaFor(&MathRequest{}),
			"response_schema": SchemaFor(&MathResponse{}),
		}))

	log.Printf("Listening on %q for %q, %q, %q, %q \n", nc.ConnectedAddr(), "math.add", "math.subtract", "math.multiply", "math.divide")

	runtime.Goexit()
}

func addHandler(req micro.Request) {
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		log.Println(err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A + calcRequest.B,
	})
}

func subtractHandler(req micro.Request) {
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		log.Println(err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A - calcRequest.B,
	})
}

func multiplyHandler(req micro.Request) {
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		log.Println(err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A * calcRequest.B,
	})
}

func divideHandler(req micro.Request) {
	var calcRequest MathRequest
	if err := json.Unmarshal(req.Data(), &calcRequest); err != nil {
		log.Println(err)
		req.Error("400", "unable to parse request", nil)
		return
	}

	req.RespondJSON(MathResponse{
		Result: calcRequest.A / calcRequest.B,
	})
}
