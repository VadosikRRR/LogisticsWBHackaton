package usecase

import (
	"context"
	"errors"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/repository"
)

type GetWindowUseCase struct {
	recordRepo repository.RawRecordRepository
}

func NewGetWindowUseCase(recordRepo repository.RawRecordRepository) *GetWindowUseCase {
	return &GetWindowUseCase{recordRepo: recordRepo}
}

func (u *GetWindowUseCase) Execute(ctx context.Context, limit int) ([]entity.RawRecord, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be positive")
	}
	return u.recordRepo.GetLatest(ctx, limit)
}
