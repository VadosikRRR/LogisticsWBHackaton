package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/artem/logisticswbhackaton/backend/platform/config"
	platformpostgres "github.com/artem/logisticswbhackaton/backend/platform/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/parquet-go/parquet-go"
)

type trainRowTime struct {
	RouteID      int64     `parquet:"route_id"`
	OfficeFromID int64     `parquet:"office_from_id"`
	Timestamp    time.Time `parquet:"timestamp"`
	Status1      int64     `parquet:"status_1"`
	Status2      int64     `parquet:"status_2"`
	Status3      int64     `parquet:"status_3"`
	Status4      int64     `parquet:"status_4"`
	Status5      int64     `parquet:"status_5"`
	Status6      int64     `parquet:"status_6"`
	Status7      int64     `parquet:"status_7"`
	Status8      int64     `parquet:"status_8"`
	Target2H     float64   `parquet:"target_2h"`
}

type trainRowInt64TS struct {
	RouteID      int64   `parquet:"route_id"`
	OfficeFromID int64   `parquet:"office_from_id"`
	Timestamp    int64   `parquet:"timestamp"`
	Status1      int64   `parquet:"status_1"`
	Status2      int64   `parquet:"status_2"`
	Status3      int64   `parquet:"status_3"`
	Status4      int64   `parquet:"status_4"`
	Status5      int64   `parquet:"status_5"`
	Status6      int64   `parquet:"status_6"`
	Status7      int64   `parquet:"status_7"`
	Status8      int64   `parquet:"status_8"`
	Target2H     float64 `parquet:"target_2h"`
}

type rawRecord struct {
	RouteID      int64
	OfficeFromID int64
	Timestamp    time.Time
	Status1      int64
	Status2      int64
	Status3      int64
	Status4      int64
	Status5      int64
	Status6      int64
	Status7      int64
	Status8      int64
	Target2H     float64
}

func main() {
	filePath := flag.String("file", defaultDatasetPath(), "path to training parquet file")
	databaseURL := flag.String("database-url", config.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:55432/logistics?sslmode=disable"), "PostgreSQL DSN")
	batchSize := flag.Int("batch-size", 5000, "rows per parquet/database batch")
	truncate := flag.Bool("truncate", false, "truncate raw_records before import")
	dryRun := flag.Bool("dry-run", false, "read and validate parquet without DB writes")
	limit := flag.Int("limit", 0, "limit total imported rows (0 = all)")
	progressEvery := flag.Int("progress-every", 50000, "print progress every N rows")
	flag.Parse()

	if *batchSize <= 0 {
		fatal("batch-size must be positive")
	}
	if *progressEvery <= 0 {
		fatal("progress-every must be positive")
	}

	fmt.Printf("importer started: file=%s batch_size=%d dry_run=%t truncate=%t limit=%d\n",
		*filePath, *batchSize, *dryRun, *truncate, *limit,
	)

	ctx := context.Background()
	var closeDB func() = func() {}
	var dbWriter dbBatchWriter = &noopWriter{}

	if !*dryRun {
		pool, err := platformpostgres.Connect(ctx, platformpostgres.Options{
			DSN:             *databaseURL,
			MaxConns:        int32(config.GetInt("POSTGRES_MAX_CONNS", 20)),
			MinConns:        int32(config.GetInt("POSTGRES_MIN_CONNS", 2)),
			ConnectTimeout:  config.GetDuration("POSTGRES_CONNECT_TIMEOUT", 5*time.Second),
			HealthcheckFreq: config.GetDuration("POSTGRES_HEALTHCHECK_PERIOD", 30*time.Second),
		})
		if err != nil {
			fatal("connect postgres: %v", err)
		}
		closeDB = pool.Close
		dbWriter = &postgresWriter{pool: pool}

		if *truncate {
			if err := truncateRawRecords(ctx, pool); err != nil {
				closeDB()
				fatal("truncate raw_records: %v", err)
			}
			fmt.Println("raw_records truncated")
		}
	}
	defer closeDB()

	imported, err := importTrainingDataset(ctx, *filePath, *batchSize, *limit, *progressEvery, dbWriter)
	if err != nil {
		fatal("import failed: %v", err)
	}

	if *dryRun {
		fmt.Printf("dry-run complete: validated_rows=%d\n", imported)
		return
	}
	fmt.Printf("import complete: inserted_rows=%d\n", imported)
}

func defaultDatasetPath() string {
	// Prefer root-relative path when launched from repository root.
	if _, err := os.Stat("datasets/train_team_track.parquet"); err == nil {
		return "datasets/train_team_track.parquet"
	}
	// Fallback for runs from backend/ directory.
	return filepath.Join("..", "datasets", "train_team_track.parquet")
}

func importTrainingDataset(
	ctx context.Context,
	filePath string,
	batchSize int,
	limit int,
	progressEvery int,
	writer dbBatchWriter,
) (int, error) {
	imported, err := importWithTimeTimestamp(ctx, filePath, batchSize, limit, progressEvery, writer)
	if err == nil {
		return imported, nil
	}
	fmt.Printf("time.Time timestamp decoding failed, fallback to int64 timestamp: %v\n", err)
	return importWithInt64Timestamp(ctx, filePath, batchSize, limit, progressEvery, writer)
}

func importWithTimeTimestamp(
	ctx context.Context,
	filePath string,
	batchSize int,
	limit int,
	progressEvery int,
	writer dbBatchWriter,
) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("open parquet file: %w", err)
	}
	defer f.Close()

	reader := parquet.NewGenericReader[trainRowTime](f)
	defer reader.Close()

	buffer := make([]trainRowTime, batchSize)
	records := make([]rawRecord, 0, batchSize)
	total := 0

	for {
		n, readErr := reader.Read(buffer)
		if n > 0 {
			records = records[:0]
			for i := 0; i < n; i++ {
				if limit > 0 && total >= limit {
					return total, nil
				}
				row := buffer[i]
				records = append(records, rawRecord{
					RouteID:      row.RouteID,
					OfficeFromID: row.OfficeFromID,
					Timestamp:    row.Timestamp.UTC(),
					Status1:      row.Status1,
					Status2:      row.Status2,
					Status3:      row.Status3,
					Status4:      row.Status4,
					Status5:      row.Status5,
					Status6:      row.Status6,
					Status7:      row.Status7,
					Status8:      row.Status8,
					Target2H:     row.Target2H,
				})
				total++
			}
			if err := writer.WriteBatch(ctx, records); err != nil {
				return total, err
			}
			if total%progressEvery < len(records) {
				fmt.Printf("progress: imported=%d\n", total)
			}
		}
		if errors.Is(readErr, io.EOF) {
			return total, nil
		}
		if readErr != nil {
			return total, readErr
		}
	}
}

func importWithInt64Timestamp(
	ctx context.Context,
	filePath string,
	batchSize int,
	limit int,
	progressEvery int,
	writer dbBatchWriter,
) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("open parquet file: %w", err)
	}
	defer f.Close()

	reader := parquet.NewGenericReader[trainRowInt64TS](f)
	defer reader.Close()

	buffer := make([]trainRowInt64TS, batchSize)
	records := make([]rawRecord, 0, batchSize)
	total := 0

	for {
		n, readErr := reader.Read(buffer)
		if n > 0 {
			records = records[:0]
			for i := 0; i < n; i++ {
				if limit > 0 && total >= limit {
					return total, nil
				}
				row := buffer[i]
				records = append(records, rawRecord{
					RouteID:      row.RouteID,
					OfficeFromID: row.OfficeFromID,
					Timestamp:    guessUnixTime(row.Timestamp),
					Status1:      row.Status1,
					Status2:      row.Status2,
					Status3:      row.Status3,
					Status4:      row.Status4,
					Status5:      row.Status5,
					Status6:      row.Status6,
					Status7:      row.Status7,
					Status8:      row.Status8,
					Target2H:     row.Target2H,
				})
				total++
			}
			if err := writer.WriteBatch(ctx, records); err != nil {
				return total, err
			}
			if total%progressEvery < len(records) {
				fmt.Printf("progress: imported=%d\n", total)
			}
		}
		if errors.Is(readErr, io.EOF) {
			return total, nil
		}
		if readErr != nil {
			return total, readErr
		}
	}
}

func guessUnixTime(v int64) time.Time {
	abs := v
	if abs < 0 {
		abs = -abs
	}
	switch {
	case abs > 1_000_000_000_000_000_000:
		return time.Unix(0, v).UTC()
	case abs > 1_000_000_000_000_000:
		return time.UnixMicro(v).UTC()
	case abs > 1_000_000_000_000:
		return time.UnixMilli(v).UTC()
	default:
		return time.Unix(v, 0).UTC()
	}
}

type dbBatchWriter interface {
	WriteBatch(ctx context.Context, records []rawRecord) error
}

type postgresWriter struct {
	pool *pgxpool.Pool
}

func (w *postgresWriter) WriteBatch(ctx context.Context, records []rawRecord) error {
	if len(records) == 0 {
		return nil
	}
	_, err := w.pool.CopyFrom(
		ctx,
		pgx.Identifier{"raw_records"},
		[]string{
			"route_id",
			"office_from_id",
			"timestamp",
			"status_1",
			"status_2",
			"status_3",
			"status_4",
			"status_5",
			"status_6",
			"status_7",
			"status_8",
			"target_2h",
		},
		pgx.CopyFromSlice(len(records), func(i int) ([]any, error) {
			record := records[i]
			return []any{
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
			}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("copy batch to raw_records: %w", err)
	}
	return nil
}

type noopWriter struct{}

func (n *noopWriter) WriteBatch(_ context.Context, _ []rawRecord) error {
	return nil
}

func truncateRawRecords(
	ctx context.Context,
	pool *pgxpool.Pool,
) error {
	if _, err := pool.Exec(ctx, "TRUNCATE TABLE raw_records"); err != nil {
		return err
	}
	return nil
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
