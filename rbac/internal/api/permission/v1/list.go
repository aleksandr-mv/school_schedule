package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (api *API) ListPermissions(ctx context.Context, req *permissionV1.ListPermissionsRequest) (*permissionV1.ListPermissionsResponse, error) {
	filter := converter.ListPermissionsToDomain(req)
	permissions, err := api.permissionService.ListPermissions(ctx, filter)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения списка прав доступа", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return converter.PermissionsToListResponse(permissions), nil
}
