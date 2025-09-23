package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	userRoleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user_role/v1"
)

func (api *API) GetRoleUsers(ctx context.Context, req *userRoleV1.GetRoleUsersRequest) (*userRoleV1.GetRoleUsersResponse, error) {
	limit := req.GetLimit()
	if limit == 0 {
		limit = 10
	}

	userIDs, nextCursor, err := api.userRoleService.GetRoleUsers(ctx, req.RoleId, limit, req.GetCursor())
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения пользователей роли", zap.Error(err))
		return nil, mapError(err)
	}

	return &userRoleV1.GetRoleUsersResponse{
		UserIds:    userIDs,
		Limit:      limit,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
	}, nil
}
