package repository

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type AggregatorClient interface {
	Ingest(ctx context.Context, batch contracts.RawRecordBatch) error
	GetStats(ctx context.Context) (int, error)
}

type DispatchClient interface {
	Run(ctx context.Context, request contracts.RunDispatchRequest) (contracts.DispatchPlan, error)
	GetLatestPlan(ctx context.Context) (contracts.DispatchPlan, bool, error)
}

type TransportClient interface {
	ListActive(ctx context.Context) ([]contracts.TransportRequest, error)
}
