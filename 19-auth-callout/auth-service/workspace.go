package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type WorkspaceUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	PhotoURL string `json:"photoURL"`
}

type WorkspaceKV struct {
	jetstream.KeyValue
}

func NewWorkspaceKV(nc *nats.Conn, bucket string) (*WorkspaceKV, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	kv, err := js.KeyValue(context.Background(), bucket)
	if err != nil {
		return nil, err
	}

	return &WorkspaceKV{kv}, nil
}

func (w *WorkspaceKV) AddUser(user *WorkspaceUser) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = w.Put(context.Background(), fmt.Sprintf("users.%s", user.Id), data)
	return err
}
