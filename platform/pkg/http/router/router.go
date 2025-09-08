// Package http_router предоставляет функции для создания и настройки HTTP роутера.
//
// Основные возможности:
// - Создание роутера на базе chi с предустановленными middleware
// - Автоматическая настройка CORS на основе конфигурации
// - Встроенные middleware для recovery, логирования, таймаутов
// - Поддержка статуса записи ответов для мониторинга
// - Автоматическое подключение health endpoints
//
// Пример использования:
//
//	router := http_router.NewRouter(httpCfg, corsCfg)
//	router.Get("/api/health", healthHandler)
//	server := http.Server{Handler: router}
package http_router

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"

	customMiddleware "github.com/aleksandr-mv/school_schedule/platform/pkg/http/middleware"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/metric"
)

// =============================================================================
// СОЗДАНИЕ И НАСТРОЙКА РОУТЕРА
// =============================================================================

// NewRouter создаёт роутер и настраивает его по конфигу: таймауты и CORS.
// Автоматически подключает стандартные middleware в оптимальном порядке:
// 1. Recovery - обработка паник
// 2. Timeout - ограничение времени выполнения запроса
// 3. HTTPMiddleware - логирование входящих запросов (включает statusRecorder)
// 4. HTTPMiddleware - сбор метрик запросов, ответов и времени выполнения
// 5. CORS - настройка политик межсайтовых запросов (если настроено)
func NewRouter(
	ctx context.Context,
	timeout time.Duration,
	allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders []string,
	allowCredentials bool,
	maxAge int,
	httpBucketBoundaries []float64,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(customMiddleware.RecoveryMiddleware)
	r.Use(customMiddleware.TimeoutMiddleware(timeout))
	r.Use(logger.HTTPMiddleware)
	r.Use(metric.HTTPMiddleware(ctx, httpBucketBoundaries))

	if len(allowedOrigins) > 0 {
		r.Use(customMiddleware.CORSMiddleware(
			allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders, allowCredentials, maxAge,
		))
	}

	return r
}
