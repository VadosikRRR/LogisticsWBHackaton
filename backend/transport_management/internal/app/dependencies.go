package app

import (
	"context"
	"fmt"
	"log/slog"

	platformpostgres "github.com/artem/logisticswbhackaton/backend/platform/postgres"
	platformredis "github.com/artem/logisticswbhackaton/backend/platform/redis"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/adapter/publisher"
	repositoryadapter "github.com/artem/logisticswbhackaton/backend/transport_management/internal/adapter/repository/postgres"
	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type serviceDependencies struct {
	requestRepo repository.TransportRequestRepository
	publisher   repository.EventPublisher
	closeFunc   func()
}

func buildDependencies(
	ctx context.Context,
	options Options,
	logger *slog.Logger,
) (serviceDependencies, error) {
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

	redisClient, err := platformredis.Connect(ctx, platformredis.Options{
		Addr:         options.RedisAddr,
		Password:     options.RedisPassword,
		DB:           options.RedisDB,
		DialTimeout:  options.RedisDialTimeout,
		ReadTimeout:  options.RedisReadTimeout,
		WriteTimeout: options.RedisWriteTimeout,
	})
	if err != nil {
		pool.Close()
		return serviceDependencies{}, err
	}

	conn, err := amqp.Dial(options.RabbitMQURL)
	if err != nil {
		pool.Close()
		_ = redisClient.Close()
		return serviceDependencies{}, fmt.Errorf("connect rabbitmq: %w", err)
	}

	rabbitPublisher, err := publisher.NewRabbitMQPublisher(logger, conn, options.RabbitMQExchange)
	if err != nil {
		pool.Close()
		_ = redisClient.Close()
		_ = conn.Close()
		return serviceDependencies{}, err
	}

	closeFunc := func() {
		_ = rabbitPublisher.Close()
		_ = redisClient.Close()
		pool.Close()
	}

	return serviceDependencies{
		requestRepo: repositoryadapter.NewTransportRequestRepository(pool, redisClient),
		publisher:   rabbitPublisher,
		closeFunc:   closeFunc,
	}, nil
}
