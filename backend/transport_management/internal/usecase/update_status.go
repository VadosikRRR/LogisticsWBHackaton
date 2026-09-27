package usecase

import (
	"context"
	"time"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/repository"
)

type UpdateStatusUseCase struct {
	requestRepo repository.TransportRequestRepository
	publisher   repository.EventPublisher
}

func NewUpdateStatusUseCase(
	requestRepo repository.TransportRequestRepository,
	publisher repository.EventPublisher,
) *UpdateStatusUseCase {
	return &UpdateStatusUseCase{
		requestRepo: requestRepo,
		publisher:   publisher,
	}
}

func (u *UpdateStatusUseCase) Execute(
	ctx context.Context,
	id string,
	status entity.Status,
) (entity.TransportRequest, error) {
	updated, err := u.requestRepo.UpdateStatus(ctx, id, status, time.Now().UTC())
	if err != nil {
		return entity.TransportRequest{}, err
	}
	if err := u.publisher.PublishStatusUpdated(ctx, updated); err != nil {
		return entity.TransportRequest{}, err
	}
	return updated, nil
}
