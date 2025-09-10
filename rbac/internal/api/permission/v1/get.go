package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (api *API) GetPermission(ctx context.Context, req *permissionV1.GetPermissionRequest) (*permissionV1.GetPermissionResponse, error) {
	permission, err := api.permissionService.GetPermission(ctx, req.PermissionId)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения права доступа", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &permissionV1.GetPermissionResponse{
		Permission: converter.PermissionToProto(permission),
	}, nil
}
