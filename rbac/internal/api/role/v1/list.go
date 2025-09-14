package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) List(ctx context.Context, req *roleV1.ListRequest) (*roleV1.ListResponse, error) {
	roles, err := api.roleService.List(ctx)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения списка ролей", zap.Error(err))
		return nil, mapError(err)
	}

	return &roleV1.ListResponse{
		Data: converter.RolesToProto(roles),
	}, nil
}
