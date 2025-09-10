package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *API) DeleteRole(ctx context.Context, req *roleV1.DeleteRoleRequest) (*emptypb.Empty, error) {
	if err := api.roleService.DeleteRole(ctx, req.RoleId); err != nil {
		logger.Error(ctx, "❌ [API] Ошибка удаления роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}
