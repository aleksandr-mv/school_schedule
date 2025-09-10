package permission

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *PermissionService) CheckPermission(ctx context.Context, userID, resource, action string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.check_permission")
	defer span.End()

	permission, err := s.permissionRepo.GetByResourceAndAction(ctx, resource, action)
	if err != nil {
		if err == model.ErrPermissionNotFound {
			span.SetStatus(codes.Error, "permission not found")
			return model.ErrPermissionDenied
		}
		logger.Error(ctx, "❌ [Service] Ошибка получения права доступа", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	roles, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения ролей пользователя", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	for _, role := range roles {
		hasPermission, err := s.rolePermissionRepo.HasPermission(ctx, role.ID.String(), permission.ID.String())
		if err != nil {
			logger.Error(ctx, "❌ [Service] Ошибка проверки права роли", zap.Error(err))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		if hasPermission {
			return nil
		}
	}

	span.SetStatus(codes.Error, "permission denied")
	return model.ErrPermissionDenied
}
