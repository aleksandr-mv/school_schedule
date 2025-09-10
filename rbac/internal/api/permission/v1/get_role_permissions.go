package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (api *API) ListPermissionsByRole(ctx context.Context, req *permissionV1.ListPermissionsByRoleRequest) (*permissionV1.ListPermissionsByRoleResponse, error) {
	value, err := converter.GetIdentifierToValue(req.Value)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка парсинга идентификатора роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	permissions, err := api.permissionService.ListPermissionsByRole(ctx, value)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения прав роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return converter.PermissionsToListByRoleResponse(permissions), nil
}
