package producer

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type producer struct {
	asyncProducer sarama.AsyncProducer
	syncProducer  sarama.SyncProducer
	topic         string
}

func NewProducer(asyncProducer sarama.AsyncProducer, topic string) *producer {
	return &producer{
		asyncProducer: asyncProducer,
		topic:         topic,
	}
}

func NewSyncProducer(syncProducer sarama.SyncProducer, topic string) *producer {
	return &producer{
		syncProducer: syncProducer,
		topic:        topic,
	}
}

func (p *producer) Send(ctx context.Context, key, value []byte) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	if p.syncProducer != nil {
		_, _, err := p.syncProducer.SendMessage(message)
		return err
	}

	if p.asyncProducer != nil {
		select {
		case p.asyncProducer.Input() <- message:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("no producer configured")
}
