package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/auth/v1"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
)

const (
	// SessionIDMetadataKey ключ для передачи ID сессии в gRPC metadata
	SessionIDMetadataKey = "session-id"
)

type contextKey string

const (
	// userContextKey ключ для хранения пользователя в контексте
	userContextKey contextKey = "user"
	// sessionIDContextKey ключ для хранения session ID в контексте
	sessionIDContextKey contextKey = "session-id"
	// userRolesWithPermissionsContextKey ключ для хранения ролей с правами пользователя в контексте
	userRolesWithPermissionsContextKey contextKey = "user-roles-with-permissions"
)

// IAMClient интерфейс для аутентификации пользователей
type IAMClient interface {
	Whoami(ctx context.Context, req *authV1.WhoamiRequest, opts ...grpc.CallOption) (*authV1.WhoamiResponse, error)
}

// AuthInterceptor interceptor для аутентификации gRPC запросов
type AuthInterceptor struct {
	iamClient IAMClient
}

// NewAuthInterceptor создает новый interceptor аутентификации
func NewAuthInterceptor(iamClient IAMClient) *AuthInterceptor {
	return &AuthInterceptor{
		iamClient: iamClient,
	}
}

// Unary возвращает unary server interceptor для аутентификации
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		authCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(authCtx, req)
	}
}

// authenticate выполняет аутентификацию и добавляет пользователя в контекст
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	// Извлекаем metadata из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Получаем session ID из metadata
	sessionIDs := md.Get(SessionIDMetadataKey)
	if len(sessionIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing session-id in metadata")
	}

	sessionID := sessionIDs[0]
	if sessionID == "" {
		return nil, status.Error(codes.Unauthenticated, "empty session-id")
	}

	// Валидируем сессию через IAM сервис
	whoamiRes, err := i.iamClient.Whoami(ctx, &authV1.WhoamiRequest{
		SessionId: sessionID,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}

	// Добавляем пользователя, session ID и роли с правами в контекст
	authCtx := context.WithValue(ctx, userContextKey, whoamiRes.User)
	authCtx = context.WithValue(authCtx, sessionIDContextKey, sessionID)
	authCtx = context.WithValue(authCtx, userRolesWithPermissionsContextKey, whoamiRes.RolesWithPermissions)

	return authCtx, nil
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	user, ok := ctx.Value(userContextKey).(*commonV1.User)
	return user, ok
}

// GetUserContextKey возвращает ключ контекста для пользователя
func GetUserContextKey() contextKey {
	return userContextKey
}

// GetUserRolesWithPermissionsContextKey возвращает ключ контекста для ролей с правами пользователя
func GetUserRolesWithPermissionsContextKey() contextKey {
	return userRolesWithPermissionsContextKey
}

// GetUserRolesWithPermissionsFromContext извлекает роли с правами пользователя из контекста
func GetUserRolesWithPermissionsFromContext(ctx context.Context) ([]*commonV1.RoleWithPermissions, bool) {
	rolesWithPermissions, ok := ctx.Value(userRolesWithPermissionsContextKey).([]*commonV1.RoleWithPermissions)
	return rolesWithPermissions, ok
}

// GetUserRolesFromContext извлекает только роли пользователя из контекста
func GetUserRolesFromContext(ctx context.Context) ([]*commonV1.Role, bool) {
	rolesWithPermissions, ok := GetUserRolesWithPermissionsFromContext(ctx)
	if !ok {
		return nil, false
	}

	roles := make([]*commonV1.Role, 0, len(rolesWithPermissions))
	for _, rwp := range rolesWithPermissions {
		roles = append(roles, rwp.Role)
	}

	return roles, true
}

// GetUserPermissionsFromContext извлекает все права пользователя из контекста
func GetUserPermissionsFromContext(ctx context.Context) ([]*commonV1.Permission, bool) {
	rolesWithPermissions, ok := GetUserRolesWithPermissionsFromContext(ctx)
	if !ok {
		return nil, false
	}

	permissionsMap := make(map[string]*commonV1.Permission)
	for _, rwp := range rolesWithPermissions {
		for _, permission := range rwp.Permissions {
			permissionsMap[permission.Id] = permission
		}
	}

	permissions := make([]*commonV1.Permission, 0, len(permissionsMap))
	for _, permission := range permissionsMap {
		permissions = append(permissions, permission)
	}

	return permissions, true
}

// GetSessionIDFromContext извлекает session ID из контекста
func GetSessionIDFromContext(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(sessionIDContextKey).(string)
	return sessionID, ok
}

// AddSessionIDToContext добавляет session ID в контекст
func AddSessionIDToContext(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDContextKey, sessionID)
}

// ForwardSessionIDToGRPC добавляет session ID из контекста в исходящие gRPC metadata
func ForwardSessionIDToGRPC(ctx context.Context) context.Context {
	sessionID, ok := GetSessionIDFromContext(ctx)
	if !ok || sessionID == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, SessionIDMetadataKey, sessionID)
}
