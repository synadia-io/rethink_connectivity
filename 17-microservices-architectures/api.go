package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"
	"github.com/segmentio/ksuid"
)

type API struct {
	ctx context.Context
	kv  jetstream.KeyValue
}

type Product struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	ReviewCount   int     `json:"review_count"`
	ReviewAverage float64 `json:"review_average"`
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
	product := Product{}
	err := json.Unmarshal(req.Data(), &product)
	if err != nil {
		req.Error("400", "Bad Request", []byte(err.Error()))
		return
	}
	id := ksuid.New().String()
	product.ID = id

	jsonProduct, err := json.Marshal(product)
	if err != nil {
		req.Error("500", "Internal Server Error", []byte(err.Error()))
		return
	}
	_, err = a.kv.Create(a.ctx, "products."+id, jsonProduct)
	if err != nil {
		req.Error("500", "Internal Server Error", []byte(err.Error()))
		return
	}

	req.RespondJSON(product)
}

func (a *API) DeleteProduct(req micro.Request) {
	id := strings.Split(req.Subject(), ".")[3]
	err := a.kv.Delete(a.ctx, "products."+id)
	if err != nil {
		req.Error("400", "Bad Request", []byte(err.Error()))
		return
	}

	req.Respond([]byte("OK"))
}

// This is a bit of a hack to get multiple entries until we
// get batch get in the next version of NATS.
func (a *API) fetchMultiple(keys string) ([]jetstream.KeyValueEntry, error) {
	results := []jetstream.KeyValueEntry{}

	watcher, err := a.kv.Watch(a.ctx, keys)
	if err != nil {
		return results, err
	}

	for entry := range watcher.Updates() {
		if entry == nil {
			break
		}
		results = append(results, entry)
	}

	return results, nil
}
