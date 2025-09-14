package role

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *RoleService) Create(ctx context.Context, name, description string) (uuid.UUID, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.create_role")
	defer span.End()

	role, err := s.roleRepo.Create(ctx, name, description)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка создания роли в репозитории", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return uuid.Nil, err
	}

	return role, nil
}
