package role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *RoleService) ListRoles(ctx context.Context, nameFilter string) ([]*model.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.list_roles")
	defer span.End()

	roles, err := s.roleRepo.List(ctx, nameFilter)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения списка ролей из репозитория", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return roles, nil
}
