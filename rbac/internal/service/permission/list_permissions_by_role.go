package permission

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *PermissionService) ListPermissionsByRole(ctx context.Context, roleValue string) ([]*model.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.list_permissions_by_role")
	defer span.End()

	permissions, err := s.rolePermissionRepo.ListPermissionsByRole(ctx, roleValue)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения прав роли", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return permissions, nil
}
