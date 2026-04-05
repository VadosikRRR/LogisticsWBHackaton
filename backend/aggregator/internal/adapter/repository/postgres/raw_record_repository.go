package postgres

import (
	"context"
	"fmt"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RawRecordRepository struct {
	pool *pgxpool.Pool
}

func NewRawRecordRepository(pool *pgxpool.Pool) *RawRecordRepository {
	return &RawRecordRepository{pool: pool}
}

func (r *RawRecordRepository) SaveBatch(ctx context.Context, records []entity.RawRecord) error {
	batch := &pgx.Batch{}
	for _, record := range records {
		batch.Queue(
			`INSERT INTO raw_records (
				route_id,
				office_from_id,
				timestamp,
				status_1,
				status_2,
				status_3,
				status_4,
				status_5,
				status_6,
				status_7,
				status_8,
				target_2h
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			record.RouteID,
			record.OfficeFromID,
			record.Timestamp,
			record.Status1,
			record.Status2,
			record.Status3,
			record.Status4,
			record.Status5,
			record.Status6,
			record.Status7,
			record.Status8,
			record.Target2H,
		)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(records); i++ {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("insert raw record %d: %w", i, err)
		}
	}
	if err := results.Close(); err != nil {
		return fmt.Errorf("close batch: %w", err)
	}
	return nil
}

func (r *RawRecordRepository) GetLatest(ctx context.Context, limit int) ([]entity.RawRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			route_id,
			office_from_id,
			timestamp,
			status_1,
			status_2,
			status_3,
			status_4,
			status_5,
			status_6,
			status_7,
			status_8,
			target_2h
		FROM raw_records
		ORDER BY timestamp DESC, id DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("query latest raw records: %w", err)
	}
	defer rows.Close()

	reversed := make([]entity.RawRecord, 0, limit)
	for rows.Next() {
		var record entity.RawRecord
		if err := rows.Scan(
			&record.RouteID,
			&record.OfficeFromID,
			&record.Timestamp,
			&record.Status1,
			&record.Status2,
			&record.Status3,
			&record.Status4,
			&record.Status5,
			&record.Status6,
			&record.Status7,
			&record.Status8,
			&record.Target2H,
		); err != nil {
			return nil, fmt.Errorf("scan raw record: %w", err)
		}
		reversed = append(reversed, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate raw records: %w", err)
	}

	// Return records in chronological order for downstream batch processing.
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	return reversed, nil
}

func (r *RawRecordRepository) Count(ctx context.Context) (int, error) {
	var count int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM raw_records`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count raw records: %w", err)
	}
	return count, nil
}
