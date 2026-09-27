package client

import (
	"context"
	"net/http"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
	"github.com/artem/logisticswbhackaton/backend/platform/httpx"
)

type AggregatorClient struct {
	baseURL    string
	httpClient *http.Client
}

type statsResponse struct {
	TotalRecords int `json:"total_records"`
}

func NewAggregatorClient(baseURL string, httpClient *http.Client) *AggregatorClient {
	return &AggregatorClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *AggregatorClient) Ingest(ctx context.Context, batch contracts.RawRecordBatch) error {
	url := httpx.JoinURL(c.baseURL, "/v1/records")
	return httpx.PostJSON(ctx, c.httpClient, url, batch, nil)
}

func (c *AggregatorClient) GetStats(ctx context.Context) (int, error) {
	url := httpx.JoinURL(c.baseURL, "/v1/stats")
	var response statsResponse
	if err := httpx.GetJSON(ctx, c.httpClient, url, &response); err != nil {
		return 0, err
	}
	return response.TotalRecords, nil
}
