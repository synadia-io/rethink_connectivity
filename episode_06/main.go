package main

import (
	"log"
	"runtime"
	"strconv"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL, nats.UserCredentials("math.creds"))
	if err != nil {
		log.Fatalln(err)
	}

	nc.QueueSubscribe("math.double", "math", func(msg *nats.Msg) {
		num, err := strconv.Atoi(string(msg.Data))
		if err != nil {
			msg.Respond([]byte("Body is not a number"))
		} else {
			msg.Respond([]byte(strconv.Itoa(num + num)))
		}
	})

	nc.QueueSubscribe("math.triple", "math", func(msg *nats.Msg) {
		num, err := strconv.Atoi(string(msg.Data))
		if err != nil {
			msg.Respond([]byte("Body is not a number"))
		} else {
			msg.Respond([]byte(strconv.Itoa(num + num + num)))
		}
	})

	nc.QueueSubscribe("math.quadruple", "math", func(msg *nats.Msg) {
		num, err := strconv.Atoi(string(msg.Data))
		if err != nil {
			msg.Respond([]byte("Body is not a number"))
		} else {
			msg.Respond([]byte(strconv.Itoa(num + num + num + num)))
		}
	})

	runtime.Goexit()
}
