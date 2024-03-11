package main

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"
)

type API struct {
	ctx context.Context
	kv  jetstream.KeyValue
}

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewAPI(ctx context.Context, kv jetstream.KeyValue) *API {
	return &API{ctx: ctx, kv: kv}
}

func (a *API) ListProducts(req micro.Request) {
	entries, err := a.fetchMultiple("products.*")
	if err != nil {
		req.Error("500", "Internal Server Error", []byte(err.Error()))
		return
	}

	products := []Product{}
	for _, entry := range entries {
		product := Product{}
		if err := json.Unmarshal(entry.Value(), &product); err != nil {
			req.Error("500", "Internal Server Error", []byte(err.Error()))
			return
		}
		products = append(products, product)
	}

	req.RespondJSON(products)
}

func (a *API) CreateProduct(req micro.Request) {
	req.Respond([]byte("Create"))
}

func (a *API) UpdateProduct(req micro.Request) {
	req.Respond([]byte("Update"))
}

func (a *API) DeleteProduct(req micro.Request) {
	req.Respond([]byte("Delete"))
}

// This is a bit of a hack to get multiple entries until we
// get batch get in the next version of NATS.
func (a *API) fetchMultiple(keys string) ([]jetstream.KeyValueEntry, error) {
	results := []jetstream.KeyValueEntry{}

	watcher, err := a.kv.Watch(a.ctx, keys)
	if err != nil {
		return results, err
	}

	for {
		select {
		case <-a.ctx.Done():
			return results, nil
		case entry := <-watcher.Updates():
			if entry == nil {
				break
			}
			results = append(results, entry)
		}
	}
}
