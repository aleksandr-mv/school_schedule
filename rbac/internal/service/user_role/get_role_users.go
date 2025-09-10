package user_role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *UserRoleService) GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor *string) ([]string, int32, *string, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_role_users")
	defer span.End()

	userIDs, totalCount, nextCursor, err := s.userRoleRepo.GetRoleUsers(ctx, roleID, limit, cursor)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения пользователей роли", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, nil, err
	}

	return userIDs, totalCount, nextCursor, nil
}
