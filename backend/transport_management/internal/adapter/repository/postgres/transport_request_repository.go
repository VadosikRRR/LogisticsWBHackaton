package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

const activeRequestsCacheKey = "transport:active_requests:v1"

type TransportRequestRepository struct {
	pool  *pgxpool.Pool
	redis *goredis.Client
}

func NewTransportRequestRepository(pool *pgxpool.Pool, redis *goredis.Client) *TransportRequestRepository {
	return &TransportRequestRepository{
		pool:  pool,
		redis: redis,
	}
}

func (r *TransportRequestRepository) UpsertFromDrafts(
	ctx context.Context,
	drafts []entity.TransportRequestDraft,
	now time.Time,
) ([]entity.TransportRequest, error) {
	if len(drafts) == 0 {
		return []entity.TransportRequest{}, nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	result := make([]entity.TransportRequest, 0, len(drafts))

	for _, draft := range drafts {
		requestID, err := newRequestID()
		if err != nil {
			return nil, err
		}
		requestKey := entity.BuildRequestKey(draft.RouteID, draft.OfficeFromID, draft.Timestamp)

		row := tx.QueryRow(ctx, `
			INSERT INTO transport_requests (
				id,
				request_key,
				route_id,
				office_from_id,
				slot_timestamp,
				required_vehicles,
				predicted_volume,
				reason,
				status,
				created_at,
				updated_at
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,'CREATED',$9,$9)
			ON CONFLICT (request_key)
			WHERE status IN ('CREATED','SENT','CONFIRMED','IN_PROGRESS')
			DO UPDATE
			SET
				required_vehicles = GREATEST(transport_requests.required_vehicles, EXCLUDED.required_vehicles),
				predicted_volume = GREATEST(transport_requests.predicted_volume, EXCLUDED.predicted_volume),
				reason = CASE WHEN EXCLUDED.reason <> '' THEN EXCLUDED.reason ELSE transport_requests.reason END,
				updated_at = EXCLUDED.updated_at
			RETURNING
				id,
				route_id,
				office_from_id,
				slot_timestamp,
				required_vehicles,
				predicted_volume,
				reason,
				status,
				created_at,
				updated_at
		`,
			requestID,
			requestKey,
			draft.RouteID,
			draft.OfficeFromID,
			draft.Timestamp,
			draft.RequiredVehicles,
			draft.PredictedVolume,
			draft.Reason,
			now,
		)

		var item entity.TransportRequest
		var status string
		if err := row.Scan(
			&item.ID,
			&item.RouteID,
			&item.OfficeFromID,
			&item.Timestamp,
			&item.RequiredVehicles,
			&item.PredictedVolume,
			&item.Reason,
			&status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("upsert transport request: %w", err)
		}
		item.Status = entity.Status(status)
		result = append(result, item)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	_ = r.invalidateActiveCache(ctx)
	return result, nil
}

func (r *TransportRequestRepository) ListActive(ctx context.Context) ([]entity.TransportRequest, error) {
	if cached, ok := r.getActiveCache(ctx); ok {
		return cached, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			id,
			route_id,
			office_from_id,
			slot_timestamp,
			required_vehicles,
			predicted_volume,
			reason,
			status,
			created_at,
			updated_at
		FROM transport_requests
		WHERE status IN ('CREATED','SENT','CONFIRMED','IN_PROGRESS')
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list active requests: %w", err)
	}
	defer rows.Close()

	items := make([]entity.TransportRequest, 0)
	for rows.Next() {
		var item entity.TransportRequest
		var status string
		if err := rows.Scan(
			&item.ID,
			&item.RouteID,
			&item.OfficeFromID,
			&item.Timestamp,
			&item.RequiredVehicles,
			&item.PredictedVolume,
			&item.Reason,
			&status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan active request: %w", err)
		}
		item.Status = entity.Status(status)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active requests: %w", err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	_ = r.setActiveCache(ctx, items)
	return items, nil
}

func (r *TransportRequestRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status entity.Status,
	now time.Time,
) (entity.TransportRequest, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE transport_requests
		SET status = $2, updated_at = $3
		WHERE id = $1
		RETURNING
			id,
			route_id,
			office_from_id,
			slot_timestamp,
			required_vehicles,
			predicted_volume,
			reason,
			status,
			created_at,
			updated_at
	`, id, string(status), now)

	var item entity.TransportRequest
	var rawStatus string
	if err := row.Scan(
		&item.ID,
		&item.RouteID,
		&item.OfficeFromID,
		&item.Timestamp,
		&item.RequiredVehicles,
		&item.PredictedVolume,
		&item.Reason,
		&rawStatus,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.TransportRequest{}, errors.New("request not found")
		}
		return entity.TransportRequest{}, fmt.Errorf("update request status: %w", err)
	}

	item.Status = entity.Status(rawStatus)
	_ = r.invalidateActiveCache(ctx)
	return item, nil
}

func (r *TransportRequestRepository) getActiveCache(ctx context.Context) ([]entity.TransportRequest, bool) {
	if r.redis == nil {
		return nil, false
	}
	payload, err := r.redis.Get(ctx, activeRequestsCacheKey).Result()
	if err != nil {
		return nil, false
	}
	var items []entity.TransportRequest
	if err := json.Unmarshal([]byte(payload), &items); err != nil {
		return nil, false
	}
	return items, true
}

func (r *TransportRequestRepository) setActiveCache(ctx context.Context, items []entity.TransportRequest) error {
	if r.redis == nil {
		return nil
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, activeRequestsCacheKey, string(raw), 30*time.Second).Err()
}

func (r *TransportRequestRepository) invalidateActiveCache(ctx context.Context) error {
	if r.redis == nil {
		return nil
	}
	return r.redis.Del(ctx, activeRequestsCacheKey).Err()
}
