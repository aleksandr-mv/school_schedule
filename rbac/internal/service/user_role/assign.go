package user_role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *UserRoleService) Assign(ctx context.Context, userID, roleID string, assignedBy *string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.assign")
	defer span.End()

	err := s.userRoleRepo.Assign(ctx, userID, roleID, assignedBy)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка назначения роли пользователю", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
