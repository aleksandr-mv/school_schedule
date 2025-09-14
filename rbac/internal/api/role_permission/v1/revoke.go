package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	rolePermissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role_permission/v1"
)

func (api *API) Revoke(ctx context.Context, req *rolePermissionV1.RevokeRequest) (*emptypb.Empty, error) {
	if err := api.rolePermissionService.Revoke(ctx, req.RoleId, req.PermissionId); err != nil {
		logger.Error(ctx, "❌ [API] Ошибка отзыва права у роли", zap.Error(err))
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
