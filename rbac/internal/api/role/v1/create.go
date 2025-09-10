package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) CreateRole(ctx context.Context, req *roleV1.CreateRoleRequest) (*roleV1.CreateRoleResponse, error) {
	roleID, err := api.roleService.CreateRole(ctx, converter.CreateRoleToDomain(req))
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка создания роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &roleV1.CreateRoleResponse{
		RoleId: roleID.String(),
	}, nil
}
