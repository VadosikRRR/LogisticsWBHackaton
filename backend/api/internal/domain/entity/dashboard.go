package entity

import (
	"time"

	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type DashboardOverview struct {
	GeneratedAt     time.Time
	TotalRawRecords int
	ActiveRequests  int
	LatestPlan      *contracts.DispatchPlan
}
