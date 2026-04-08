package service

import (
	"math"
	"time"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
)

type Planner struct{}

func NewPlanner() *Planner {
	return &Planner{}
}

func (p *Planner) BuildRequests(
	predictions []entity.Prediction,
	activeRequests []entity.ActiveTransportRequest,
	policies entity.CapacityPolicySet,
) []entity.PlannedTransportRequest {
	activeVehiclesByKey := make(map[string]int, len(activeRequests))
	for _, request := range activeRequests {
		key := entity.BuildActiveRequestKey(request.RouteID, request.OfficeFromID, request.Timestamp)
		activeVehiclesByKey[key] += request.RequiredVehicles
	}

	requests := make([]entity.PlannedTransportRequest, 0, len(predictions))

	for _, prediction := range predictions {
		policy := policies.ForOffice(prediction.OfficeFromID)
		requiredVehicles := requiredVehicles(prediction.Target2H, policy)
		if requiredVehicles == 0 {
			continue
		}

		key := entity.BuildPredictionKey(prediction.RouteID, prediction.OfficeFromID, prediction.Timestamp)
		activeVehicles := activeVehiclesByKey[key]
		needed := requiredVehicles - activeVehicles
		if needed <= 0 {
			continue
		}

		requests = append(requests, entity.PlannedTransportRequest{
			RouteID:          prediction.RouteID,
			OfficeFromID:     prediction.OfficeFromID,
			Timestamp:        prediction.Timestamp.Truncate(30 * time.Minute),
			RequiredVehicles: needed,
			PredictedVolume:  prediction.Target2H,
			Reason:           "AUTO_DISPATCH_FROM_FORECAST",
		})
	}

	return requests
}

func requiredVehicles(predictedVolume float64, policy entity.OfficeCapacityPolicy) int {
	if predictedVolume <= 0 || policy.VehicleCapacity <= 0 {
		return 0
	}

	volumeWithBuffer := predictedVolume * (1 + policy.SafetyBufferFraction)
	vehicles := int(math.Ceil(volumeWithBuffer / policy.VehicleCapacity))
	if vehicles < 1 {
		vehicles = 1
	}
	if policy.MaxVehiclesPerSlot > 0 && vehicles > policy.MaxVehiclesPerSlot {
		vehicles = policy.MaxVehiclesPerSlot
	}
	return vehicles
}
