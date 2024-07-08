package main

import (
	"testing"

	"github.com/nats-io/nats.go"
	"time"
)

func BenchmarkRequestReplyLoopback(b *testing.B) {
	nc, ns, err := RunEmbeddedServer(false, false)
	if err != nil {
		b.Fatal(err)
	}

	nc.Subscribe("hello.world", func(msg *nats.Msg) {
		msg.Respond([]byte("Hi there"))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := nc.Request("hello.world", []byte("hihi"), 10*time.Second)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	ns.Shutdown()
	ns.WaitForShutdown()
}

func BenchmarkRequestReplyInProcess(b *testing.B) {
	nc, ns, err := RunEmbeddedServer(true, false)
	if err != nil {
		b.Fatal(err)
	}

	nc.Subscribe("hello.world", func(msg *nats.Msg) {
		msg.Respond([]byte("Hi there"))
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := nc.Request("hello.world", []byte("hihi"), 10*time.Second)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	ns.Shutdown()
	ns.WaitForShutdown()
}
