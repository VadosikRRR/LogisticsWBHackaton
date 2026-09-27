package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	ConnectTimeout  time.Duration
	HealthcheckFreq time.Duration
}

func Connect(ctx context.Context, options Options) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(options.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	if options.MaxConns > 0 {
		cfg.MaxConns = options.MaxConns
	}
	if options.MinConns > 0 {
		cfg.MinConns = options.MinConns
	}
	if options.HealthcheckFreq > 0 {
		cfg.HealthCheckPeriod = options.HealthcheckFreq
	}
	if options.ConnectTimeout > 0 {
		cfg.ConnConfig.ConnectTimeout = options.ConnectTimeout
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("new postgres pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}
