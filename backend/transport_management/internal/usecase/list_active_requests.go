package usecase

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/repository"
)

type ListActiveRequestsUseCase struct {
	requestRepo repository.TransportRequestRepository
}

func NewListActiveRequestsUseCase(requestRepo repository.TransportRequestRepository) *ListActiveRequestsUseCase {
	return &ListActiveRequestsUseCase{requestRepo: requestRepo}
}

func (u *ListActiveRequestsUseCase) Execute(ctx context.Context) ([]entity.TransportRequest, error) {
	return u.requestRepo.ListActive(ctx)
}
