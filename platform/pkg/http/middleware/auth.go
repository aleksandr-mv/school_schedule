package middleware

import (
	"context"
	"net/http"

	grpcAuth "github.com/aleksandr-mv/school_schedule/platform/pkg/grpc/interceptor"
	authV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
)

const SessionIDHeader = "X-Session-Id"

// IAMClient интерфейс для аутентификации пользователей
type IAMClient interface {
	Whoami(ctx context.Context, sessionID string) (*authV1.WhoamiResponse, error)
}

// AuthMiddleware middleware для аутентификации HTTP запросов
type AuthMiddleware struct {
	iamClient IAMClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(iamClient IAMClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

// client (X-Session-Uuid) -> auth middleware (add session_uuid in ctx (incomming)) -> order api (outgoing)-> auth interceptor ->inventory
// Handle обрабатывает HTTP запрос с аутентификацией
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем session ID из заголовка
		sessionID := r.Header.Get(SessionIDHeader)
		if sessionID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.iamClient.Whoami(r.Context(), sessionID)
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session ID в контекст используя функции из grpc middleware
		ctx := r.Context()
		ctx = grpcAuth.AddSessionIDToContext(ctx, sessionID)

		// Также добавляем пользователя и права в контекст
		ctx = context.WithValue(ctx, grpcAuth.GetUserContextKey(), whoamiRes.User)
		ctx = context.WithValue(ctx, grpcAuth.GetPermissionsContextKey(), whoamiRes.Permissions)

		// Передаем управление следующему handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

// GetSessionIDFromContext извлекает session ID из контекста
func GetSessionIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionIDFromContext(ctx)
}
