package role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *RoleService) GetRole(ctx context.Context, value string) (*model.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_role")
	defer span.End()

	role, err := s.roleRepo.Get(ctx, value)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения роли из репозитория", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return role, nil
}
