package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (api *API) AssignPermissionToRole(ctx context.Context, req *permissionV1.AssignPermissionToRoleRequest) (*emptypb.Empty, error) {
	err := api.permissionService.AssignPermissionToRole(ctx, req.RoleId, req.PermissionId)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка назначения права роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}
