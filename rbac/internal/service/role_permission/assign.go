package role_permission

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *RolePermissionService) Assign(ctx context.Context, roleID, permissionID string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.assign_permission_to_role")
	defer span.End()

	err := s.rolePermissionRepo.Assign(ctx, roleID, permissionID)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка назначения права роли", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
