package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	roleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role/v1"
)

func (api *API) Delete(ctx context.Context, req *roleV1.DeleteRequest) (*emptypb.Empty, error) {
	if err := api.roleService.Delete(ctx, req.RoleId); err != nil {
		logger.Error(ctx, "❌ [API] Ошибка удаления роли", zap.Error(err))
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
