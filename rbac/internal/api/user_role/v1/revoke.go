package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	userRoleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user_role/v1"
)

func (api *API) Revoke(ctx context.Context, req *userRoleV1.RevokeRequest) (*emptypb.Empty, error) {
	if err := api.userRoleService.Revoke(ctx, req.UserId, req.RoleId); err != nil {
		logger.Error(ctx, "❌ [API] Ошибка отзыва роли у пользователя", zap.Error(err))
		return nil, mapError(err)
	}

	return &emptypb.Empty{}, nil
}
