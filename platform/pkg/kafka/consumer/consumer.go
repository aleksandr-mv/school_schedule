package consumer

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/kafka"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type consumer struct {
	group       sarama.ConsumerGroup
	topics      []string
	logger      kafka.Logger
	middlewares []Middleware
}

// NewConsumer — создаёт новый consumer.
func NewConsumer(group sarama.ConsumerGroup, topics []string, logger kafka.Logger, middlewares ...Middleware) *consumer {
	return &consumer{
		group:       group,
		topics:      topics,
		logger:      logger,
		middlewares: middlewares,
	}
}

// Consume запускает консьюмер для списка топиков.
func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	newGroupHandler := NewGroupHandler(handler, c.logger, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error(ctx, "Kafka consume error", zap.Error(err))
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}

func (c *consumer) Close() error {
	return c.group.Close()
}
