package user_role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *UserRoleService) RevokeRole(ctx context.Context, userID, roleID string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.revoke_role")
	defer span.End()

	err := s.userRoleRepo.RevokeRole(ctx, userID, roleID)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка отзыва роли у пользователя", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
