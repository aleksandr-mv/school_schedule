package permission

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *PermissionService) GetPermission(ctx context.Context, value string) (*model.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_permission")
	defer span.End()

	permission, err := s.permissionRepo.Get(ctx, value)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения права доступа из репозитория", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return permission, nil
}
