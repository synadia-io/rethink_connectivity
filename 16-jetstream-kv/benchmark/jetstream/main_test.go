package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func BenchmarkJetStreamKeyValuePutDistinct(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			err := put(fmt.Sprintf("key.%d", rand.Int()), "hello world")
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkRedisSetDistinct(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			err := set(fmt.Sprintf("key:%d", rand.Int()), "hello world")
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
