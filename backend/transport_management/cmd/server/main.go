package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/artem/logisticswbhackaton/backend/platform/config"
	"github.com/artem/logisticswbhackaton/backend/platform/logging"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/app"
)

func main() {
	logger := logging.New("transport-management")
	options := app.Options{
		Port:                    config.GetInt("TRANSPORT_PORT", 8083),
		PostgresDSN:             config.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/logistics?sslmode=disable"),
		PostgresMaxConns:        int32(config.GetInt("POSTGRES_MAX_CONNS", 20)),
		PostgresMinConns:        int32(config.GetInt("POSTGRES_MIN_CONNS", 2)),
		PostgresConnectTimeout:  config.GetDuration("POSTGRES_CONNECT_TIMEOUT", 5*time.Second),
		PostgresHealthcheckFreq: config.GetDuration("POSTGRES_HEALTHCHECK_PERIOD", 30*time.Second),
		RedisAddr:               config.GetString("REDIS_ADDR", "localhost:6379"),
		RedisPassword:           config.GetString("REDIS_PASSWORD", ""),
		RedisDB:                 config.GetInt("REDIS_DB", 0),
		RedisDialTimeout:        config.GetDuration("REDIS_DIAL_TIMEOUT", 3*time.Second),
		RedisReadTimeout:        config.GetDuration("REDIS_READ_TIMEOUT", 2*time.Second),
		RedisWriteTimeout:       config.GetDuration("REDIS_WRITE_TIMEOUT", 2*time.Second),
		RabbitMQURL:             config.GetString("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		RabbitMQExchange:        config.GetString("RABBITMQ_EXCHANGE", "transport.events"),
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
