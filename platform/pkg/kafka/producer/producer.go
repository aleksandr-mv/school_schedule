package producer

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/kafka"
)

type producer struct {
	syncProducer sarama.SyncProducer
	topic        string
	logger       kafka.Logger
}

func NewProducer(syncProducer sarama.SyncProducer, topic string, logger kafka.Logger) *producer {
	return &producer{
		syncProducer: syncProducer,
		topic:        topic,
		logger:       logger,
	}
}

func (p *producer) Send(ctx context.Context, key, value []byte) error {
	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		p.logger.Error(ctx, "Failed to send message", zap.Error(err))
		return err
	}

	p.logger.Info(ctx, "Message sent",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.Int("key_length", len(key)),
		zap.Int("value_length", len(value)),
	)

	return nil
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
