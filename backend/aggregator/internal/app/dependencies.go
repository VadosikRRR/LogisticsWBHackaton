package app

import (
	"context"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/adapter/repository/postgres"
	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/domain/repository"
	platformpostgres "github.com/artem/logisticswbhackaton/backend/platform/postgres"
)

type serviceDependencies struct {
	recordRepo repository.RawRecordRepository
	closeFunc  func()
}

func usecaseDependencies(ctx context.Context, options Options) (serviceDependencies, error) {
	pool, err := platformpostgres.Connect(ctx, platformpostgres.Options{
		DSN:             options.PostgresDSN,
		MaxConns:        options.PostgresMaxConns,
		MinConns:        options.PostgresMinConns,
		ConnectTimeout:  options.PostgresConnectTimeout,
		HealthcheckFreq: options.PostgresHealthcheckFreq,
	})
	if err != nil {
		return serviceDependencies{}, err
	}

	return serviceDependencies{
		recordRepo: postgres.NewRawRecordRepository(pool),
		closeFunc:  pool.Close,
	}, nil
}
