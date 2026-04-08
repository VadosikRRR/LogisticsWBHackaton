package usecase

import (
	"context"
	"time"

	"github.com/artem/logisticswbhackaton/backend/api/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/api/internal/domain/repository"
)

type GetDashboardOverviewUseCase struct {
	aggregatorClient repository.AggregatorClient
	dispatchClient   repository.DispatchClient
	transportClient  repository.TransportClient
}

func NewGetDashboardOverviewUseCase(
	aggregatorClient repository.AggregatorClient,
	dispatchClient repository.DispatchClient,
	transportClient repository.TransportClient,
) *GetDashboardOverviewUseCase {
	return &GetDashboardOverviewUseCase{
		aggregatorClient: aggregatorClient,
		dispatchClient:   dispatchClient,
		transportClient:  transportClient,
	}
}

func (u *GetDashboardOverviewUseCase) Execute(ctx context.Context) (entity.DashboardOverview, error) {
	totalRawRecords, err := u.aggregatorClient.GetStats(ctx)
	if err != nil {
		return entity.DashboardOverview{}, err
	}

	activeRequests, err := u.transportClient.ListActive(ctx)
	if err != nil {
		return entity.DashboardOverview{}, err
	}

	latestPlan, found, err := u.dispatchClient.GetLatestPlan(ctx)
	if err != nil {
		return entity.DashboardOverview{}, err
	}

	overview := entity.DashboardOverview{
		GeneratedAt:     time.Now().UTC(),
		TotalRawRecords: totalRawRecords,
		ActiveRequests:  len(activeRequests),
	}
	if found {
		overview.LatestPlan = &latestPlan
	}
	return overview, nil
}
