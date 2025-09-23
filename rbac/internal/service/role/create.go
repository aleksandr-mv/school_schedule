package role

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/errreport"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/google/uuid"
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
