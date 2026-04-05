package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/artem/logisticswbhackaton/backend/transport_management/internal/domain/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultExchangeType = "topic"

type RabbitMQPublisher struct {
	logger   *slog.Logger
	conn     *amqp.Connection
	exchange string
}

type Event struct {
	EventType  string                  `json:"event_type"`
	OccurredAt time.Time               `json:"occurred_at"`
	Request    entity.TransportRequest `json:"request"`
}

func NewRabbitMQPublisher(
	logger *slog.Logger,
	conn *amqp.Connection,
	exchange string,
) (*RabbitMQPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		exchange,
		defaultExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("declare exchange: %w", err)
	}

	return &RabbitMQPublisher{
		logger:   logger,
		conn:     conn,
		exchange: exchange,
	}, nil
}

func (p *RabbitMQPublisher) PublishCreated(ctx context.Context, requests []entity.TransportRequest) error {
	for _, request := range requests {
		if err := p.publish(ctx, "transport.request.created", Event{
			EventType:  "TRANSPORT_REQUEST_CREATED",
			OccurredAt: time.Now().UTC(),
			Request:    request,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p *RabbitMQPublisher) PublishStatusUpdated(ctx context.Context, request entity.TransportRequest) error {
	return p.publish(ctx, "transport.request.updated", Event{
		EventType:  "TRANSPORT_REQUEST_STATUS_UPDATED",
		OccurredAt: time.Now().UTC(),
		Request:    request,
	})
}

func (p *RabbitMQPublisher) Close() error {
	if p.conn == nil {
		return nil
	}
	return p.conn.Close()
}

func (p *RabbitMQPublisher) publish(ctx context.Context, routingKey string, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("open rabbitmq channel: %w", err)
	}
	defer ch.Close()

	if err := ch.PublishWithContext(
		ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			Timestamp:   time.Now().UTC(),
		},
	); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	p.logger.Info("event published", "routing_key", routingKey, "request_id", event.Request.ID)
	return nil
}
