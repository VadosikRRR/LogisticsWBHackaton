package dto

import (
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type RunDispatchRequest = contracts.RunDispatchRequest
type RunDispatchResponse = contracts.RunDispatchResponse

type LatestPlanResponse struct {
	Plan contracts.DispatchPlan `json:"plan"`
}

func FromEntityPlan(plan entity.DispatchPlan) contracts.DispatchPlan {
	requests := make([]contracts.CreateTransportRequest, 0, len(plan.Requests))
	for _, request := range plan.Requests {
		requests = append(requests, contracts.CreateTransportRequest{
			RouteID:          request.RouteID,
			OfficeFromID:     request.OfficeFromID,
			Timestamp:        request.Timestamp,
			RequiredVehicles: request.RequiredVehicles,
			PredictedVolume:  request.PredictedVolume,
			Reason:           request.Reason,
		})
	}

	return contracts.DispatchPlan{
		GeneratedAt:      plan.GeneratedAt,
		WindowLimit:      plan.WindowLimit,
		RecordsScanned:   plan.RecordsScanned,
		PredictionsCount: plan.PredictionsCount,
		ModelVersion:     plan.ModelVersion,
		CreatedRequests:  plan.CreatedRequests,
		Requests:         requests,
	}
}
