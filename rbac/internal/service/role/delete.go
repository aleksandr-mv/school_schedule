package role

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/errreport"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *RoleService) Delete(ctx context.Context, id string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.delete_role")
	defer span.End()

	err := s.roleRepo.Delete(ctx, id)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка удаления роли из репозитория", err)
		return err
	}

	return nil
}
