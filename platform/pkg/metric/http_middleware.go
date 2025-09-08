package metric

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

// HTTPMiddleware создает HTTP middleware для сбора метрик.
// Middleware собирает метрики запросов, ответов и времени выполнения.
// bucketBoundaries - настраиваемые границы для гистограммы времени ответа
func HTTPMiddleware(ctx context.Context, bucketBoundaries []float64) func(http.Handler) http.Handler {
	// Проверяем, инициализирован ли MeterProvider
	meterProvider := otel.GetMeterProvider()
	if meterProvider == nil {
		// Метрики не инициализированы - возвращаем простой middleware без метрик
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	// Создаем инструменты метрик один раз при инициализации middleware
	meter := otel.Meter("http-server")

	requestCounter, err := meter.Int64Counter(
		getMetricName("http_requests_total"),
		metric.WithDescription("Количество HTTP запросов"),
	)
	if err != nil {
		globalMetrics.logger.Error(ctx, "❌ [Metrics] Ошибка создания request counter", zap.Error(err))
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	responseCounter, err := meter.Int64Counter(
		getMetricName("http_responses_total"),
		metric.WithDescription("Количество HTTP ответов"),
	)
	if err != nil {
		globalMetrics.logger.Error(ctx, "❌ [Metrics] Ошибка создания response counter", zap.Error(err))
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	responseTimeHistogram, err := meter.Float64Histogram(
		getMetricName("http_response_time_seconds"),
		metric.WithDescription("Время выполнения HTTP запросов"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(bucketBoundaries...),
	)
	if err != nil {
		globalMetrics.logger.Error(ctx, "❌ [Metrics] Ошибка создания response time histogram", zap.Error(err))
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Увеличиваем счетчик входящих запросов
			requestCounter.Add(r.Context(), 1)

			// Засекаем время начала обработки
			startTime := time.Now()

			// Создаем wrapper для ResponseWriter для отслеживания статуса
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

			// Выполняем обработчик
			next.ServeHTTP(sw, r)
			duration := time.Since(startTime)

			// Определяем статус ответа
			status := "success"
			if sw.status >= 400 {
				status = "error"
			}

			// Увеличиваем счетчик ответов с атрибутами
			responseCounter.Add(r.Context(), 1,
				metric.WithAttributes(
					attribute.String("status", status),
					attribute.String("method", r.Method),
					attribute.String("path", r.URL.Path),
				),
			)

			// Записываем время выполнения в гистограмму
			responseTimeHistogram.Record(r.Context(), duration.Seconds(),
				metric.WithAttributes(
					attribute.String("status", status),
				),
			)
		})
	}
}

// statusWriter оборачивает http.ResponseWriter для отслеживания статуса ответа
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}
