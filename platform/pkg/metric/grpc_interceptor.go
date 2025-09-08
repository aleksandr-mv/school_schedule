package metric

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor создает gRPC unary interceptor для сбора метрик.
// Interceptor собирает метрики запросов, ответов и времени выполнения.
// bucketBoundaries - настраиваемые границы для гистограммы времени ответа
func UnaryServerInterceptor(ctx context.Context, bucketBoundaries []float64) grpc.UnaryServerInterceptor {
	// Создаем инструменты метрик один раз при инициализации интерсептора
	meter := GetMeterProvider().Meter("grpc-server")

	requestCounter, err := meter.Int64Counter(
		getMetricName("grpc_requests_total"),
		metric.WithDescription("Количество gRPC запросов"),
	)
	if err != nil {
		globalMetrics.logger.Error(ctx, "❌ [Metrics] Ошибка создания gRPC request counter", zap.Error(err))
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	responseCounter, err := meter.Int64Counter(
		getMetricName("grpc_responses_total"),
		metric.WithDescription("Количество gRPC ответов"),
	)
	if err != nil {
		globalMetrics.logger.Error(ctx, "❌ [Metrics] Ошибка создания gRPC response counter", zap.Error(err))
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	responseTimeHistogram, err := meter.Float64Histogram(
		getMetricName("grpc_response_time_seconds"),
		metric.WithDescription("Время выполнения gRPC запросов"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(bucketBoundaries...),
	)
	if err != nil {
		globalMetrics.logger.Error(ctx, "❌ [Metrics] Ошибка создания gRPC response time histogram", zap.Error(err))
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Увеличиваем счетчик входящих запросов
		requestCounter.Add(ctx, 1)

		// Засекаем время начала обработки
		startTime := time.Now()

		// Выполняем обработчик
		resp, err := handler(ctx, req)
		duration := time.Since(startTime)

		// Определяем статус ответа
		status := "success"
		if err != nil {
			status = "error"
		}

		// Увеличиваем счетчик ответов с атрибутами
		responseCounter.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("status", status),
				attribute.String("method", info.FullMethod),
			),
		)

		// Записываем время выполнения в гистограмму
		responseTimeHistogram.Record(ctx, duration.Seconds(),
			metric.WithAttributes(
				attribute.String("status", status),
			),
		)

		return resp, err
	}
}
