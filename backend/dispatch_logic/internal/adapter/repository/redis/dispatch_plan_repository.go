package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
	goredis "github.com/redis/go-redis/v9"
)

const latestDispatchPlanKey = "dispatch:latest_plan:v1"

type DispatchPlanRepository struct {
	client *goredis.Client
}

func NewDispatchPlanRepository(client *goredis.Client) *DispatchPlanRepository {
	return &DispatchPlanRepository{client: client}
}

func (r *DispatchPlanRepository) SaveLatest(ctx context.Context, plan entity.DispatchPlan) error {
	raw, err := json.Marshal(plan)
	if err != nil {
		return fmt.Errorf("marshal dispatch plan: %w", err)
	}
	if err := r.client.Set(ctx, latestDispatchPlanKey, raw, 0).Err(); err != nil {
		return fmt.Errorf("save latest plan to redis: %w", err)
	}
	return nil
}

func (r *DispatchPlanRepository) GetLatest(ctx context.Context) (entity.DispatchPlan, bool, error) {
	raw, err := r.client.Get(ctx, latestDispatchPlanKey).Bytes()
	if err == goredis.Nil {
		return entity.DispatchPlan{}, false, nil
	}
	if err != nil {
		return entity.DispatchPlan{}, false, fmt.Errorf("get latest plan from redis: %w", err)
	}

	var plan entity.DispatchPlan
	if err := json.Unmarshal(raw, &plan); err != nil {
		return entity.DispatchPlan{}, false, fmt.Errorf("unmarshal dispatch plan: %w", err)
	}
	return plan, true, nil
}
