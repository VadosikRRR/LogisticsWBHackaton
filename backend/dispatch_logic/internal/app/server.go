package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/adapter/client"
	repositoryadapter "github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/adapter/repository/redis"
	httpcontroller "github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/controller/http"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/entity"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/repository"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/domain/service"
	"github.com/artem/logisticswbhackaton/backend/dispatch_logic/internal/usecase"
	platformredis "github.com/artem/logisticswbhackaton/backend/platform/redis"
)

type Server struct {
	httpServer *http.Server
	closeFunc  func()
}

func NewDefaultServer(ctx context.Context, options Options, logger *slog.Logger) (*Server, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	aggregatorClient := client.NewAggregatorClient(options.AggregatorURL, httpClient)
	var mlClient repository.MLClient
	if options.UseMockML {
		mlClient = client.NewMockMLClient()
	} else {
		mlClient = client.NewMLClient(options.MLBackendURL, httpClient)
	}
	transportClient := client.NewTransportClient(options.TransportURL, httpClient)

	redisClient, err := platformredis.Connect(ctx, platformredis.Options{
		Addr:         options.RedisAddr,
		Password:     options.RedisPassword,
		DB:           options.RedisDB,
		DialTimeout:  options.RedisDialTimeout,
		ReadTimeout:  options.RedisReadTimeout,
		WriteTimeout: options.RedisWriteTimeout,
	})
	if err != nil {
		return nil, err
	}

	planRepo := repositoryadapter.NewDispatchPlanRepository(redisClient)
	planner := service.NewPlanner()

	policies := entity.CapacityPolicySet{
		DefaultPolicy: entity.OfficeCapacityPolicy{
			VehicleCapacity:      options.VehicleCapacity,
			SafetyBufferFraction: options.SafetyBuffer,
			MaxVehiclesPerSlot:   options.MaxVehiclesPerSlot,
		},
	}

	runUC := usecase.NewRunDispatchCycleUseCase(
		aggregatorClient,
		mlClient,
		transportClient,
		planRepo,
		planner,
		policies,
		options.DefaultWindowLimit,
	)
	latestUC := usecase.NewGetLatestPlanUseCase(planRepo)

	handler := httpcontroller.NewHandler(runUC, latestUC, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", options.Port),
		Handler:           loggingMiddleware(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		closeFunc: func() {
			_ = redisClient.Close()
		},
	}, nil
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
