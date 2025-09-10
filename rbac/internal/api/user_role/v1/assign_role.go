package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (api *API) AssignRole(ctx context.Context, req *userRoleV1.AssignRoleRequest) (*emptypb.Empty, error) {
	err := api.userRoleService.AssignRole(ctx, req.UserId, req.RoleId, req.AssignedBy)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка назначения роли пользователю", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}
