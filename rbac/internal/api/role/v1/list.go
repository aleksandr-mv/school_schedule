package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) ListRoles(ctx context.Context, req *roleV1.ListRolesRequest) (*roleV1.ListRolesResponse, error) {
	nameFilter := converter.ParseNameFilter(req)
	roles, err := api.roleService.ListRoles(ctx, nameFilter)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения списка ролей", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &roleV1.ListRolesResponse{
		Data: converter.RolesToProto(roles),
	}, nil
}
