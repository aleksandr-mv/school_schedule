package role

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *RoleService) Delete(ctx context.Context, id string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.delete_role")
	defer span.End()

	err := s.roleRepo.Delete(ctx, id)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка удаления роли из репозитория", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
