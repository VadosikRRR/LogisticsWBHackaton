package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
	"github.com/artem/logisticswbhackaton/backend/platform/httpx"
)

type AggregatorClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAggregatorClient(baseURL string, httpClient *http.Client) *AggregatorClient {
	return &AggregatorClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *AggregatorClient) GetWindow(ctx context.Context, limit int) ([]contracts.RawRecord, error) {
	url := httpx.JoinURL(c.baseURL, fmt.Sprintf("/v1/records/window?limit=%d", limit))

	var response contracts.RawRecordBatch
	if err := httpx.GetJSON(ctx, c.httpClient, url, &response); err != nil {
		return nil, err
	}
	return response.Records, nil
}
