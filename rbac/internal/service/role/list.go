package role

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *RoleService) List(ctx context.Context) ([]*model.Role, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.list_roles")
	defer span.End()

	roles, err := s.roleRepo.List(ctx)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения списка ролей из репозитория", err)
		return nil, err
	}

	return roles, nil
}
