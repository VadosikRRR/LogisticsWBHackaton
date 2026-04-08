package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/app"
	"github.com/artem/logisticswbhackaton/backend/platform/config"
	"github.com/artem/logisticswbhackaton/backend/platform/logging"
)

func main() {
	logger := logging.New("dispatch-logic")
	options := app.Options{
		Port:               config.GetInt("DISPATCH_PORT", 8082),
		AggregatorURL:      config.GetString("AGGREGATOR_URL", "http://localhost:8081"),
		MLBackendURL:       config.GetString("ML_BACKEND_URL", "http://localhost:8000"),
		TransportURL:       config.GetString("TRANSPORT_URL", "http://localhost:8083"),
		UseMockML:          config.GetBool("DISPATCH_USE_MOCK_ML", false),
		RedisAddr:          config.GetString("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      config.GetString("REDIS_PASSWORD", ""),
		RedisDB:            config.GetInt("REDIS_DB", 0),
		RedisDialTimeout:   config.GetDuration("REDIS_DIAL_TIMEOUT", 3*time.Second),
		RedisReadTimeout:   config.GetDuration("REDIS_READ_TIMEOUT", 2*time.Second),
		RedisWriteTimeout:  config.GetDuration("REDIS_WRITE_TIMEOUT", 2*time.Second),
		DefaultWindowLimit: config.GetInt("DISPATCH_DEFAULT_WINDOW_LIMIT", 14000),
		VehicleCapacity:    config.GetFloat("DISPATCH_VEHICLE_CAPACITY", 20.0),
		SafetyBuffer:       config.GetFloat("DISPATCH_SAFETY_BUFFER", 0.15),
		MaxVehiclesPerSlot: config.GetInt("DISPATCH_MAX_VEHICLES_PER_SLOT", 100),
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
