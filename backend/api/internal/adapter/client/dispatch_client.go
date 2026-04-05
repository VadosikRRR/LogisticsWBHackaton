package client

import (
	"context"
	"net/http"
	"strings"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
	"github.com/artem/logisticswbhackaton/backend/platform/httpx"
)

type DispatchClient struct {
	baseURL    string
	httpClient *http.Client
}

type latestPlanResponse struct {
	Plan contracts.DispatchPlan `json:"plan"`
}

func NewDispatchClient(baseURL string, httpClient *http.Client) *DispatchClient {
	return &DispatchClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *DispatchClient) Run(
	ctx context.Context,
	request contracts.RunDispatchRequest,
) (contracts.DispatchPlan, error) {
	url := httpx.JoinURL(c.baseURL, "/v1/dispatch/run")
	var response contracts.RunDispatchResponse
	if err := httpx.PostJSON(ctx, c.httpClient, url, request, &response); err != nil {
		return contracts.DispatchPlan{}, err
	}
	return response.Plan, nil
}

func (c *DispatchClient) GetLatestPlan(ctx context.Context) (contracts.DispatchPlan, bool, error) {
	url := httpx.JoinURL(c.baseURL, "/v1/dispatch/plans/latest")
	var response latestPlanResponse
	if err := httpx.GetJSON(ctx, c.httpClient, url, &response); err != nil {
		if strings.Contains(err.Error(), "unexpected status 404") {
			return contracts.DispatchPlan{}, false, nil
		}
		return contracts.DispatchPlan{}, false, err
	}
	return response.Plan, true, nil
}
