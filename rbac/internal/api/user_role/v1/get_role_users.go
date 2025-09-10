package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (api *API) GetRoleUsers(ctx context.Context, req *userRoleV1.GetRoleUsersRequest) (*userRoleV1.GetRoleUsersResponse, error) {
	limit := int32(10)
	if req.Limit != nil {
		limit = *req.Limit
	}

	userIDs, totalCount, nextCursor, err := api.userRoleService.GetRoleUsers(ctx, req.RoleId, limit, req.Cursor)
	if err != nil {
		logger.Error(ctx, "❌ [API] Ошибка получения пользователей роли", zap.Error(err))
		return nil, mapError(ctx, err)
	}

	return &userRoleV1.GetRoleUsersResponse{
		UserIds:    userIDs,
		TotalCount: totalCount,
		Limit:      limit,
		NextCursor: nextCursor,
		HasMore:    nextCursor != nil,
	}, nil
}
