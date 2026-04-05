package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/repository"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/service"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type RunDispatchCycleUseCase struct {
	aggregatorClient repository.AggregatorClient
	mlClient         repository.MLClient
	transportClient  repository.TransportClient
	planRepo         repository.DispatchPlanRepository
	planner          *service.Planner
	policies         entity.CapacityPolicySet
	defaultLimit     int
}

func NewRunDispatchCycleUseCase(
	aggregatorClient repository.AggregatorClient,
	mlClient repository.MLClient,
	transportClient repository.TransportClient,
	planRepo repository.DispatchPlanRepository,
	planner *service.Planner,
	policies entity.CapacityPolicySet,
	defaultLimit int,
) *RunDispatchCycleUseCase {
	return &RunDispatchCycleUseCase{
		aggregatorClient: aggregatorClient,
		mlClient:         mlClient,
		transportClient:  transportClient,
		planRepo:         planRepo,
		planner:          planner,
		policies:         policies,
		defaultLimit:     defaultLimit,
	}
}

func (u *RunDispatchCycleUseCase) Execute(ctx context.Context, requestedLimit int) (entity.DispatchPlan, error) {
	limit := requestedLimit
	if limit <= 0 {
		limit = u.defaultLimit
	}
	if limit <= 0 {
		return entity.DispatchPlan{}, errors.New("window limit must be positive")
	}

	rawRecords, err := u.aggregatorClient.GetWindow(ctx, limit)
	if err != nil {
		return entity.DispatchPlan{}, err
	}

	if len(rawRecords) == 0 {
		plan := entity.DispatchPlan{
			GeneratedAt:      time.Now().UTC(),
			WindowLimit:      limit,
			RecordsScanned:   0,
			PredictionsCount: 0,
			CreatedRequests:  0,
			Requests:         []entity.PlannedTransportRequest{},
		}
		if err := u.planRepo.SaveLatest(ctx, plan); err != nil {
			return entity.DispatchPlan{}, err
		}
		return plan, nil
	}

	predictResponse, err := u.mlClient.Predict(ctx, rawRecords)
	if err != nil {
		return entity.DispatchPlan{}, err
	}

	predictions := toEntityPredictions(predictResponse.Predictions)
	activeTransport, err := u.transportClient.ListActive(ctx)
	if err != nil {
		return entity.DispatchPlan{}, err
	}

	plannedRequests := u.planner.BuildRequests(predictions, toEntityActiveRequests(activeTransport), u.policies)
	if len(plannedRequests) > 0 {
		_, err = u.transportClient.CreateBulk(ctx, contracts.CreateTransportRequestsRequest{
			Requests: toCreateTransportRequests(plannedRequests),
		})
		if err != nil {
			return entity.DispatchPlan{}, err
		}
	}

	plan := entity.DispatchPlan{
		GeneratedAt:      time.Now().UTC(),
		WindowLimit:      limit,
		RecordsScanned:   len(rawRecords),
		PredictionsCount: len(predictions),
		ModelVersion:     predictResponse.ModelVersion,
		CreatedRequests:  len(plannedRequests),
		Requests:         plannedRequests,
	}

	if err := u.planRepo.SaveLatest(ctx, plan); err != nil {
		return entity.DispatchPlan{}, err
	}
	return plan, nil
}

func toEntityPredictions(predictions []contracts.Prediction) []entity.Prediction {
	out := make([]entity.Prediction, 0, len(predictions))
	for _, prediction := range predictions {
		out = append(out, entity.Prediction{
			RouteID:      prediction.RouteID,
			OfficeFromID: prediction.OfficeFromID,
			Timestamp:    prediction.Timestamp,
			Target2H:     prediction.Target2H,
		})
	}
	return out
}

func toEntityActiveRequests(requests []contracts.TransportRequest) []entity.ActiveTransportRequest {
	out := make([]entity.ActiveTransportRequest, 0, len(requests))
	for _, request := range requests {
		out = append(out, entity.ActiveTransportRequest{
			RouteID:          request.RouteID,
			OfficeFromID:     request.OfficeFromID,
			Timestamp:        request.Timestamp,
			RequiredVehicles: request.RequiredVehicles,
		})
	}
	return out
}

func toCreateTransportRequests(requests []entity.PlannedTransportRequest) []contracts.CreateTransportRequest {
	out := make([]contracts.CreateTransportRequest, 0, len(requests))
	for _, request := range requests {
		out = append(out, contracts.CreateTransportRequest{
			RouteID:          request.RouteID,
			OfficeFromID:     request.OfficeFromID,
			Timestamp:        request.Timestamp,
			RequiredVehicles: request.RequiredVehicles,
			PredictedVolume:  request.PredictedVolume,
			Reason:           request.Reason,
		})
	}
	return out
}
