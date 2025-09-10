package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) GetRole(ctx context.Context, req *roleV1.GetRoleRequest) (*roleV1.GetRoleResponse, error) {
	roleValue, err := converter.GetIdentifierToValue(req.Value)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка парсинга идентификатора роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	role, err := api.roleService.GetRole(ctx, roleValue)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &roleV1.GetRoleResponse{
		Role: converter.RoleToProto(role),
	}, nil
}
