package role

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *RoleService) Update(ctx context.Context, updateRole *model.UpdateRole) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.update_role")
	defer span.End()

	err := s.roleRepo.Update(ctx, updateRole)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка обновления роли в репозитории", err)
		return err
	}

	return nil
}
