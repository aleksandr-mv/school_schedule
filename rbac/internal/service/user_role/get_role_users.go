package user_role

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
)

func (s *UserRoleService) GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor string) ([]string, *string, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_role_users")
	defer span.End()

	userIDs, nextCursor, err := s.userRoleRepo.GetRoleUsers(ctx, roleID, limit, cursor)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения пользователей роли", err)
		return nil, nil, err
	}

	return userIDs, nextCursor, nil
}
