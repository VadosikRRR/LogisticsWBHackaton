package inmemory

import (
	"context"
	"sync"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
)

type DispatchPlanRepository struct {
	mu     sync.RWMutex
	latest entity.DispatchPlan
	found  bool
}

func NewDispatchPlanRepository() *DispatchPlanRepository {
	return &DispatchPlanRepository{}
}

func (r *DispatchPlanRepository) SaveLatest(_ context.Context, plan entity.DispatchPlan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.latest = plan
	r.found = true
	return nil
}

func (r *DispatchPlanRepository) GetLatest(_ context.Context) (entity.DispatchPlan, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.latest, r.found, nil
}
