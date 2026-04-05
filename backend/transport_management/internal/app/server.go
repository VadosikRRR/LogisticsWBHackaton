package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	httpcontroller "github.com/artem/logisticswbhackaton/backend/transport_management/internal/controller/http"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/usecase"
)

type Server struct {
	httpServer *http.Server
	closeFunc  func()
}

func NewDefaultServer(ctx context.Context, options Options, logger *slog.Logger) (*Server, error) {
	deps, err := buildDependencies(ctx, options, logger)
	if err != nil {
		return nil, err
	}

	createUC := usecase.NewCreateBulkRequestsUseCase(deps.requestRepo, deps.publisher)
	activeUC := usecase.NewListActiveRequestsUseCase(deps.requestRepo)
	updateUC := usecase.NewUpdateStatusUseCase(deps.requestRepo, deps.publisher)

	handler := httpcontroller.NewHandler(createUC, activeUC, updateUC, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", options.Port),
		Handler:           loggingMiddleware(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{httpServer: httpServer, closeFunc: deps.closeFunc}, nil
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	if s.closeFunc != nil {
		s.closeFunc()
	}
	return nil
}
