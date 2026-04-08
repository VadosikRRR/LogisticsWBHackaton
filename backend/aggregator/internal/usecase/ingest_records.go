package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/repository"
)

type IngestRecordsUseCase struct {
	recordRepo repository.RawRecordRepository
}

func NewIngestRecordsUseCase(recordRepo repository.RawRecordRepository) *IngestRecordsUseCase {
	return &IngestRecordsUseCase{recordRepo: recordRepo}
}

func (u *IngestRecordsUseCase) Execute(ctx context.Context, records []entity.RawRecord) error {
	if len(records) == 0 {
		return errors.New("records batch is empty")
	}

	for i, record := range records {
		if err := validateRawRecord(record); err != nil {
			return fmt.Errorf("record %d: %w", i, err)
		}
	}

	return u.recordRepo.SaveBatch(ctx, records)
}

func validateRawRecord(record entity.RawRecord) error {
	if record.RouteID <= 0 {
		return errors.New("route_id must be positive")
	}
	if record.OfficeFromID <= 0 {
		return errors.New("office_from_id must be positive")
	}
	if record.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if record.Timestamp.After(time.Now().Add(365 * 24 * time.Hour)) {
		return errors.New("timestamp too far in future")
	}
	if math.IsNaN(record.Target2H) || math.IsInf(record.Target2H, 0) {
		return errors.New("target_2h must be finite")
	}
	if record.Target2H < 0 {
		return errors.New("target_2h must be non-negative")
	}
	return nil
}
