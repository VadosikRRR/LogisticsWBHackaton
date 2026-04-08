package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/artem/logisticswbhackaton/backend/api/internal/adapter/client"
	httpcontroller "github.com/artem/logisticswbhackaton/backend/api/internal/controller/http"
	"github.com/artem/logisticswbhackaton/backend/api/internal/usecase"
)

type Server struct {
	httpServer *http.Server
}

func NewDefaultServer(options Options, logger *slog.Logger) *Server {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	aggregatorClient := client.NewAggregatorClient(options.AggregatorURL, httpClient)
	dispatchClient := client.NewDispatchClient(options.DispatchURL, httpClient)
	transportClient := client.NewTransportClient(options.TransportURL, httpClient)

	ingestUC := usecase.NewIngestRecordsUseCase(aggregatorClient)
	triggerDispatchUC := usecase.NewTriggerDispatchUseCase(dispatchClient)
	dashboardUC := usecase.NewGetDashboardOverviewUseCase(aggregatorClient, dispatchClient, transportClient)

	handler := httpcontroller.NewHandler(ingestUC, triggerDispatchUC, dashboardUC, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", options.Port),
		Handler:           loggingMiddleware(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{httpServer: httpServer}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
