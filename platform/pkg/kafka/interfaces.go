package kafka

import (
	"context"

	"go.uber.org/zap"
)

// Logger интерфейс для логирования
type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
}

// Metrics интерфейс для метрик
type Metrics interface {
	IncrementCounter(name string, tags map[string]string)
	IncrementGauge(name string, tags map[string]string, value float64)
}
