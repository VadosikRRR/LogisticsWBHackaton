package usecase

import (
	"context"
	"errors"

	"github.com/artem/logisticswbhackaton/backend/api/internal/domain/repository"
	"github.com/artem/logisticswbhackaton/backend/platform/contracts"
)

type IngestRecordsUseCase struct {
	aggregatorClient repository.AggregatorClient
}

func NewIngestRecordsUseCase(aggregatorClient repository.AggregatorClient) *IngestRecordsUseCase {
	return &IngestRecordsUseCase{aggregatorClient: aggregatorClient}
}

func (u *IngestRecordsUseCase) Execute(ctx context.Context, batch contracts.RawRecordBatch) error {
	if len(batch.Records) == 0 {
		return errors.New("records batch is empty")
	}
	return u.aggregatorClient.Ingest(ctx, batch)
}
