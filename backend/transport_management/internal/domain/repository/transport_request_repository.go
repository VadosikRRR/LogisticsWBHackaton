package repository

import (
	"context"
	"time"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
)

type TransportRequestRepository interface {
	UpsertFromDrafts(ctx context.Context, drafts []entity.TransportRequestDraft, now time.Time) ([]entity.TransportRequest, error)
	ListActive(ctx context.Context) ([]entity.TransportRequest, error)
	UpdateStatus(ctx context.Context, id string, status entity.Status, now time.Time) (entity.TransportRequest, error)
}
