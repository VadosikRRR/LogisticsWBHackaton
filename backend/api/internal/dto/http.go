package dto

import (
	"time"

	"github.com/artem/logisticswbhackaton/backend/api/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type IngestWarehouseRequest = contracts.RawRecordBatch
type RunDispatchRequest = contracts.RunDispatchRequest

type DashboardOverviewResponse struct {
	GeneratedAt     time.Time               `json:"generated_at"`
	TotalRawRecords int                     `json:"total_raw_records"`
	ActiveRequests  int                     `json:"active_requests"`
	LatestPlan      *contracts.DispatchPlan `json:"latest_plan,omitempty"`
}

func FromEntityDashboard(overview entity.DashboardOverview) DashboardOverviewResponse {
	return DashboardOverviewResponse{
		GeneratedAt:     overview.GeneratedAt,
		TotalRawRecords: overview.TotalRawRecords,
		ActiveRequests:  overview.ActiveRequests,
		LatestPlan:      overview.LatestPlan,
	}
}
