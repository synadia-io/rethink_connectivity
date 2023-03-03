package main

import (
	"encoding/json"

	"github.com/cdipaolo/sentiment"
	"github.com/invopop/jsonschema"
	"github.com/nats-io/nats.go/micro"
)

var model sentiment.Models

type SentimentRequest struct {
	Text string `json:"text"`
}

type SentimentResponse struct {
	*sentiment.Analysis
}

// Schema returns a nats micro compatible json schema for sentiment request and sentiment response
func Schema() (*micro.Schema, error) {
	reqSchema, err := json.Marshal(jsonschema.Reflect(&SentimentRequest{}))
	if err != nil {
		return nil, err
	}

	resSchema, err := json.Marshal(jsonschema.Reflect(&SentimentResponse{}))
	if err != nil {
		return nil, err
	}

	return &micro.Schema{
		Request:  string(reqSchema),
		Response: string(resSchema),
	}, nil
}

func AnalyzeSentiment(req *SentimentRequest, model sentiment.Models) *SentimentResponse {
	a := model.SentimentAnalysis(req.Text, sentiment.English)
	return &SentimentResponse{Analysis: a}
}
