package usecase

import (
	"context"
	"errors"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/repository"
)

type GetLatestPlanUseCase struct {
	planRepo repository.DispatchPlanRepository
}

func NewGetLatestPlanUseCase(planRepo repository.DispatchPlanRepository) *GetLatestPlanUseCase {
	return &GetLatestPlanUseCase{planRepo: planRepo}
}

func (u *GetLatestPlanUseCase) Execute(ctx context.Context) (entity.DispatchPlan, error) {
	plan, found, err := u.planRepo.GetLatest(ctx)
	if err != nil {
		return entity.DispatchPlan{}, err
	}
	if !found {
		return entity.DispatchPlan{}, errors.New("dispatch plan not found")
	}
	return plan, nil
}
