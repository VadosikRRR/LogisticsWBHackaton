package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	httpcontroller "github.com/artem/logisticswbhackaton/backend/aggregator/internal/controller/http"
	"github.com/artem/logisticswbhackaton/backend/aggregator/internal/usecase"
)

type Server struct {
	httpServer *http.Server
	closeFunc  func()
}

func NewServer(port int, handler *httpcontroller.Handler, logger *slog.Logger, closeFunc func()) *Server {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           loggingMiddleware(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{httpServer: httpServer, closeFunc: closeFunc}
}

func NewDefaultServer(ctx context.Context, options Options, logger *slog.Logger) (*Server, error) {
	deps, err := usecaseDependencies(ctx, options)
	if err != nil {
		return nil, err
	}

	ingestUC := usecase.NewIngestRecordsUseCase(deps.recordRepo)
	windowUC := usecase.NewGetWindowUseCase(deps.recordRepo)
	statsUC := usecase.NewGetStatsUseCase(deps.recordRepo)

	handler := httpcontroller.NewHandler(ingestUC, windowUC, statsUC, logger)
	return NewServer(options.Port, handler, logger, deps.closeFunc), nil
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
