package repository

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
)

type DispatchPlanRepository interface {
	SaveLatest(ctx context.Context, plan entity.DispatchPlan) error
	GetLatest(ctx context.Context) (entity.DispatchPlan, bool, error)
}
