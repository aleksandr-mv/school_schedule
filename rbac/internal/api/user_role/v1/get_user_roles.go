package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (api *API) GetUserRoles(ctx context.Context, req *userRoleV1.GetUserRolesRequest) (*userRoleV1.GetUserRolesResponse, error) {
	roles, err := api.userRoleService.GetUserRoles(ctx, req.UserId)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения ролей пользователя", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &userRoleV1.GetUserRolesResponse{
		Data: converter.RolesToProto(roles),
	}, nil
}
