package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/app"
	"github.com/artem/logisticswbhackaton/backend/platform/config"
	"github.com/artem/logisticswbhackaton/backend/platform/logging"
)

func main() {
	logger := logging.New("data-aggregator")
	options := app.Options{
		Port:                    config.GetInt("AGGREGATOR_PORT", 8081),
		PostgresDSN:             config.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/logistics?sslmode=disable"),
		PostgresMaxConns:        int32(config.GetInt("POSTGRES_MAX_CONNS", 20)),
		PostgresMinConns:        int32(config.GetInt("POSTGRES_MIN_CONNS", 2)),
		PostgresConnectTimeout:  config.GetDuration("POSTGRES_CONNECT_TIMEOUT", 5*time.Second),
		PostgresHealthcheckFreq: config.GetDuration("POSTGRES_HEALTHCHECK_PERIOD", 30*time.Second),
	}

	server, err := app.NewDefaultServer(context.Background(), options, logger)
	if err != nil {
		logger.Error("failed to init server", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server started", "port", options.Port)
		errCh <- server.Run()
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
