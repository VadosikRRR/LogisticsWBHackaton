package publisher

import (
	"context"
	"log/slog"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
)

type LogPublisher struct {
	logger *slog.Logger
}

func NewLogPublisher(logger *slog.Logger) *LogPublisher {
	return &LogPublisher{logger: logger}
}

func (p *LogPublisher) PublishCreated(_ context.Context, requests []entity.TransportRequest) error {
	p.logger.Info("transport requests created/updated", "count", len(requests))
	return nil
}

func (p *LogPublisher) PublishStatusUpdated(_ context.Context, request entity.TransportRequest) error {
	p.logger.Info("transport request status updated", "id", request.ID, "status", request.Status)
	return nil
}
