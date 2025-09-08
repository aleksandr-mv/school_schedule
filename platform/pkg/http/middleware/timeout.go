package middleware

import (
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware устанавливает таймаут выполнения обработчика (не путать с socket I/O таймаутами сервера).
// Если timeout <= 0, таймаут не применяется.
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if timeout <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
