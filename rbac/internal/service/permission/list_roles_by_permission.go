package permission

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *PermissionService) ListRolesByPermission(ctx context.Context, permissionValue string) ([]*model.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.list_roles_by_permission")
	defer span.End()

	roles, err := s.rolePermissionRepo.ListRolesByPermission(ctx, permissionValue)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения ролей по праву", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return roles, nil
}
