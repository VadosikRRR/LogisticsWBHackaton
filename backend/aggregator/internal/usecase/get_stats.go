package usecase

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/repository"
)

type GetStatsUseCase struct {
	recordRepo repository.RawRecordRepository
}

func NewGetStatsUseCase(recordRepo repository.RawRecordRepository) *GetStatsUseCase {
	return &GetStatsUseCase{recordRepo: recordRepo}
}

func (u *GetStatsUseCase) Execute(ctx context.Context) (int, error) {
	return u.recordRepo.Count(ctx)
}
