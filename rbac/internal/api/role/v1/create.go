package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	roleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) Create(ctx context.Context, req *roleV1.CreateRequest) (*roleV1.CreateResponse, error) {
	id, err := api.roleService.Create(ctx, req.GetName(), req.GetDescription())
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка создания роли", zap.Error(err))
		return nil, mapError(err)
	}

	return &roleV1.CreateResponse{
		RoleId: id.String(),
	}, nil
}
