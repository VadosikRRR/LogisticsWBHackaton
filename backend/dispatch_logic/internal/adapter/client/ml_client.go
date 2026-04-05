package client

import (
	"context"
	"net/http"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
	"github.com/artem/logisticswbhackaton/backend/platform/httpx"
)

type MLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewMLClient(baseURL string, httpClient *http.Client) *MLClient {
	return &MLClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *MLClient) Predict(ctx context.Context, records []contracts.RawRecord) (contracts.PredictResponse, error) {
	url := httpx.JoinURL(c.baseURL, "/v1/predict")
	request := contracts.PredictRequest{
		Records: records,
	}

	var response contracts.PredictResponse
	if err := httpx.PostJSON(ctx, c.httpClient, url, request, &response); err != nil {
		return contracts.PredictResponse{}, err
	}
	return response, nil
}
