package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/repository"
)

type CreateBulkRequestsUseCase struct {
	requestRepo repository.TransportRequestRepository
	publisher   repository.EventPublisher
}

func NewCreateBulkRequestsUseCase(
	requestRepo repository.TransportRequestRepository,
	publisher repository.EventPublisher,
) *CreateBulkRequestsUseCase {
	return &CreateBulkRequestsUseCase{
		requestRepo: requestRepo,
		publisher:   publisher,
	}
}

func (u *CreateBulkRequestsUseCase) Execute(
	ctx context.Context,
	drafts []entity.TransportRequestDraft,
) ([]entity.TransportRequest, error) {
	if len(drafts) == 0 {
		return []entity.TransportRequest{}, nil
	}

	for i, draft := range drafts {
		if err := draft.Validate(); err != nil {
			return nil, fmt.Errorf("draft %d: %w", i, err)
		}
	}

	created, err := u.requestRepo.UpsertFromDrafts(ctx, drafts, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	if err := u.publisher.PublishCreated(ctx, created); err != nil {
		return nil, err
	}
	return created, nil
}
