package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (api *API) RevokePermissionFromRole(ctx context.Context, req *permissionV1.RevokePermissionFromRoleRequest) (*emptypb.Empty, error) {
	if err := api.permissionService.RevokePermissionFromRole(ctx, req.RoleId, req.PermissionId); err != nil {
		logger.Error(ctx, "❌ [API] Ошибка отзыва права у роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}
