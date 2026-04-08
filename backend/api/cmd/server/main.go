package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/artem/logisticswbhackaton/backend/api/internal/app"
	"github.com/artem/logisticswbhackaton/backend/platform/config"
	"github.com/artem/logisticswbhackaton/backend/platform/logging"
)

func main() {
	logger := logging.New("api-gateway")
	options := app.Options{
		Port:          config.GetInt("API_PORT", 8080),
		AggregatorURL: config.GetString("AGGREGATOR_URL", "http://localhost:8081"),
		DispatchURL:   config.GetString("DISPATCH_URL", "http://localhost:8082"),
		TransportURL:  config.GetString("TRANSPORT_URL", "http://localhost:8083"),
	}

	server := app.NewDefaultServer(options, logger)

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
