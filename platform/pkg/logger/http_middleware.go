package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HTTPMiddleware создает HTTP middleware для логирования запросов.
// Middleware логирует начало/конец HTTP-запроса с минимально необходимыми полями.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		Info(ctx, "🚀 [HTTP] Запрос начат",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		startTime := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sr, r)
		duration := time.Since(startTime)

		if sr.status >= 400 {
			Error(ctx, "❌ [HTTP] Запрос завершился ошибкой",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sr.status),
				zap.Duration("duration", duration),
			)
		} else {
			Info(ctx, "✅ [HTTP] Запрос завершён",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", sr.status),
				zap.Duration("duration", duration),
			)
		}
	})
}

// statusRecorder оборачивает http.ResponseWriter для отслеживания статуса и размера ответа
type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	n, err := sr.ResponseWriter.Write(b)
	sr.bytes += n
	return n, err
}
