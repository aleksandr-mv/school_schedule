package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/converter"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) Get(ctx context.Context, req *roleV1.GetRequest) (*roleV1.GetResponse, error) {
	enrichedRole, err := api.roleService.Get(ctx, req.RoleId)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения роли", zap.Error(err))
		return nil, mapError(err)
	}

	return &roleV1.GetResponse{
		Data: converter.EnrichedRoleToProto(enrichedRole),
	}, nil
}
