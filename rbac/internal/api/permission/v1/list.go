package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (api *API) List(ctx context.Context, req *permissionV1.ListRequest) (*permissionV1.ListResponse, error) {
	permissions, err := api.permissionService.List(ctx)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения списка прав доступа", zap.Error(err))
		return nil, mapError(err)
	}

	return &permissionV1.ListResponse{
		Data: converter.PermissionsToProto(permissions),
	}, nil
}
