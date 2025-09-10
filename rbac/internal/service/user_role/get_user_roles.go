package user_role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *UserRoleService) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_user_roles")
	defer span.End()

	roles, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения ролей пользователя", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return roles, nil
}
