package main

import (
	"github.com/invopop/jsonschema"
)

type MathRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type MathResponse struct {
	Result int `json:"result"`
}

func SchemaFor(t any) string {
	schema := jsonschema.Reflect(t)
	data, _ := schema.MarshalJSON()
	return string(data)
}
