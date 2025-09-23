package middleware

import (
	"context"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

// RecoveryMiddleware перехватывает паники в HTTP-обработчиках, логирует их и отдаёт 500 Internal Server Error.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer func(c context.Context) {
			if rec := recover(); rec != nil {
				logger.Error(c, "💥 [HTTP] Panic при обработке запроса",
					zap.Any("panic", rec),
					zap.ByteString("stacktrace", debug.Stack()),
				)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}(ctx)

		next.ServeHTTP(w, r)
	})
}
