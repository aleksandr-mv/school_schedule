package role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *RoleService) UpdateRole(ctx context.Context, updateRole *model.UpdateRole) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.update_role")
	defer span.End()

	if err := updateRole.Validate(); err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка валидации роли", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return model.ErrInvalidCredentials
	}

	err := s.roleRepo.Update(ctx, updateRole)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка обновления роли в репозитории", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
