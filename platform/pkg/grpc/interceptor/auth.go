package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// Заголовки от Envoy External Auth (после успешной аутентификации)
	// Минимальный набор: session ID для идентификации и permissions для авторизации
	HeaderSessionID       = "x-session-id"
	HeaderUserPermissions = "x-user-permissions"

	// HTTP заголовки
	HeaderCookie        = "cookie"
	HeaderAuthorization = "authorization"
	HeaderContentType   = "content-type"
	HeaderAuthStatus    = "X-Auth-Status"

	// Cookies
	SessionCookieName = "X-Session-Id"

	// Значения
	ContentTypeJSON  = "application/json"
	AuthStatusDenied = "denied"
)

type contextKey string

const (
	// sessionIDContextKey ключ для хранения session ID в контексте
	sessionIDContextKey contextKey = "session-id"
	// userPermissionsStringsContextKey ключ для хранения прав как строк
	userPermissionsStringsContextKey contextKey = "user-permissions-strings"
)

// AuthInterceptor interceptor для чтения данных пользователя из Envoy заголовков
// Работает ТОЛЬКО с защищенными методами (публичные уже отфильтрованы PublicFilter)
type AuthInterceptor struct{}

// NewAuthInterceptor создает новый interceptor аутентификации
func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

// Unary возвращает unary server interceptor для аутентификации
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if IsPublicMethod(ctx) {
			return handler(ctx, req)
		}

		authCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(authCtx, req)
	}
}

// authenticate читает данные пользователя из Envoy заголовков и создает контекст
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	sessionIDs := md.Get(HeaderSessionID)
	if len(sessionIDs) == 0 || sessionIDs[0] == "" {
		return nil, status.Error(codes.Unauthenticated, "missing session ID from Envoy")
	}

	var permissions []string
	if permHeaders := md.Get(HeaderUserPermissions); len(permHeaders) > 0 && permHeaders[0] != "" {
		permissions = strings.Split(permHeaders[0], ",")
	}

	authCtx := context.WithValue(ctx, sessionIDContextKey, sessionIDs[0])
	authCtx = context.WithValue(authCtx, userPermissionsStringsContextKey, permissions)

	return authCtx, nil
}

// GetSessionIDContextKey возвращает ключ для session ID в контексте
func GetSessionIDContextKey() contextKey {
	return sessionIDContextKey
}

// GetSessionIDFromContext извлекает session ID из контекста
func GetSessionIDFromContext(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(sessionIDContextKey).(string)
	return sessionID, ok
}

// GetUserPermissionsStringsFromContext извлекает права как строки из контекста для авторизации
func GetUserPermissionsStringsFromContext(ctx context.Context) ([]string, bool) {
	permissions, ok := ctx.Value(userPermissionsStringsContextKey).([]string)
	return permissions, ok
}
