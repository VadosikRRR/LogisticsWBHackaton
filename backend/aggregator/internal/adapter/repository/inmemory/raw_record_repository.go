package inmemory

import (
	"context"
	"sync"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/entity"
)

type RawRecordRepository struct {
	mu      sync.RWMutex
	records []entity.RawRecord
}

func NewRawRecordRepository() *RawRecordRepository {
	return &RawRecordRepository{
		records: make([]entity.RawRecord, 0, 1024),
	}
}

func (r *RawRecordRepository) SaveBatch(_ context.Context, records []entity.RawRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = append(r.records, records...)
	return nil
}

func (r *RawRecordRepository) GetLatest(_ context.Context, limit int) ([]entity.RawRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit > len(r.records) {
		limit = len(r.records)
	}

	start := len(r.records) - limit
	out := make([]entity.RawRecord, limit)
	copy(out, r.records[start:])
	return out, nil
}

func (r *RawRecordRepository) Count(_ context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.records), nil
}
