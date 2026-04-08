package usecase

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/api/internal/domain/repository"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type TriggerDispatchUseCase struct {
	dispatchClient repository.DispatchClient
}

func NewTriggerDispatchUseCase(dispatchClient repository.DispatchClient) *TriggerDispatchUseCase {
	return &TriggerDispatchUseCase{dispatchClient: dispatchClient}
}

func (u *TriggerDispatchUseCase) Execute(
	ctx context.Context,
	request contracts.RunDispatchRequest,
) (contracts.DispatchPlan, error) {
	return u.dispatchClient.Run(ctx, request)
}
