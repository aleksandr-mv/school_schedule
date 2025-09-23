package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	rolePermissionV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role_permission/v1"
)

func (api *API) Assign(ctx context.Context, req *rolePermissionV1.AssignRequest) (*emptypb.Empty, error) {
	err := api.rolePermissionService.Assign(ctx, req.RoleId, req.PermissionId)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка назначения права роли", zap.Error(err))
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
