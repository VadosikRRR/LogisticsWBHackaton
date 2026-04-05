package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
)

type TransportRequestRepository struct {
	mu          sync.RWMutex
	sequence    uint64
	requests    map[string]entity.TransportRequest
	activeIndex map[string]string
}

func NewTransportRequestRepository() *TransportRequestRepository {
	return &TransportRequestRepository{
		requests:    make(map[string]entity.TransportRequest),
		activeIndex: make(map[string]string),
	}
}

func (r *TransportRequestRepository) UpsertFromDrafts(
	_ context.Context,
	drafts []entity.TransportRequestDraft,
	now time.Time,
) ([]entity.TransportRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]entity.TransportRequest, 0, len(drafts))

	for _, draft := range drafts {
		key := entity.BuildRequestKey(draft.RouteID, draft.OfficeFromID, draft.Timestamp)
		if id, found := r.activeIndex[key]; found {
			existing := r.requests[id]
			if draft.RequiredVehicles > existing.RequiredVehicles {
				existing.RequiredVehicles = draft.RequiredVehicles
			}
			if draft.PredictedVolume > existing.PredictedVolume {
				existing.PredictedVolume = draft.PredictedVolume
			}
			if draft.Reason != "" {
				existing.Reason = draft.Reason
			}
			existing.UpdatedAt = now
			r.requests[id] = existing
			result = append(result, existing)
			continue
		}

		r.sequence++
		id := fmt.Sprintf("tr-%d", r.sequence)
		created := entity.TransportRequest{
			ID:               id,
			RouteID:          draft.RouteID,
			OfficeFromID:     draft.OfficeFromID,
			Timestamp:        draft.Timestamp,
			RequiredVehicles: draft.RequiredVehicles,
			PredictedVolume:  draft.PredictedVolume,
			Reason:           draft.Reason,
			Status:           entity.StatusCreated,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		r.requests[id] = created
		r.activeIndex[key] = id
		result = append(result, created)
	}

	return result, nil
}

func (r *TransportRequestRepository) ListActive(_ context.Context) ([]entity.TransportRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]entity.TransportRequest, 0, len(r.requests))
	for _, request := range r.requests {
		if request.IsActive() {
			out = append(out, request)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})

	return out, nil
}

func (r *TransportRequestRepository) UpdateStatus(
	_ context.Context,
	id string,
	status entity.Status,
	now time.Time,
) (entity.TransportRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	request, found := r.requests[id]
	if !found {
		return entity.TransportRequest{}, errors.New("request not found")
	}

	request.Status = status
	request.UpdatedAt = now
	r.requests[id] = request

	key := entity.BuildRequestKey(request.RouteID, request.OfficeFromID, request.Timestamp)
	if request.IsActive() {
		r.activeIndex[key] = id
	} else {
		delete(r.activeIndex, key)
	}
	return request, nil
}
