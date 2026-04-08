package repository

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type AggregatorClient interface {
	GetWindow(ctx context.Context, limit int) ([]contracts.RawRecord, error)
}

type MLClient interface {
	Predict(ctx context.Context, records []contracts.RawRecord) (contracts.PredictResponse, error)
}

type TransportClient interface {
	ListActive(ctx context.Context) ([]contracts.TransportRequest, error)
	CreateBulk(ctx context.Context, request contracts.CreateTransportRequestsRequest) ([]contracts.TransportRequest, error)
}
