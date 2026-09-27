package client

import (
	"context"
	"net/http"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
	"github.com/artem/logisticswbhackaton/backend/platform/httpx"
)

type TransportClient struct {
	baseURL    string
	httpClient *http.Client
}

type listActiveResponse struct {
	Requests []contracts.TransportRequest `json:"requests"`
}

func NewTransportClient(baseURL string, httpClient *http.Client) *TransportClient {
	return &TransportClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *TransportClient) ListActive(ctx context.Context) ([]contracts.TransportRequest, error) {
	url := httpx.JoinURL(c.baseURL, "/v1/requests/active")
	var response listActiveResponse
	if err := httpx.GetJSON(ctx, c.httpClient, url, &response); err != nil {
		return nil, err
	}
	return response.Requests, nil
}
