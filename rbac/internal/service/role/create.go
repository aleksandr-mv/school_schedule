package role

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
)

func (s *RoleService) Create(ctx context.Context, name, description string) (uuid.UUID, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.create_role")
	defer span.End()

	role, err := s.roleRepo.Create(ctx, name, description)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка создания роли в репозитории", err)
		return uuid.Nil, err
	}

	return role, nil
}
