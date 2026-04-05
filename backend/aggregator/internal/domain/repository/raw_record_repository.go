package repository

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/entity"
)

type RawRecordRepository interface {
	SaveBatch(ctx context.Context, records []entity.RawRecord) error
	GetLatest(ctx context.Context, limit int) ([]entity.RawRecord, error)
	Count(ctx context.Context) (int, error)
}
