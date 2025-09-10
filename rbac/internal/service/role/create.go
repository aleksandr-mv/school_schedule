package role

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *RoleService) CreateRole(ctx context.Context, createRole *model.CreateRole) (uuid.UUID, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.create_role")
	defer span.End()

	if err := createRole.Validate(); err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка валидации роли", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return uuid.Nil, model.ErrInvalidCredentials
	}

	role, err := s.roleRepo.Create(ctx, createRole)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка создания роли в репозитории", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return uuid.Nil, err
	}

	return role.ID, nil
}
