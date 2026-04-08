package repository

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
)

type EventPublisher interface {
	PublishCreated(ctx context.Context, requests []entity.TransportRequest) error
	PublishStatusUpdated(ctx context.Context, request entity.TransportRequest) error
}
