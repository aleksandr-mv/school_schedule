package user_role

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/errreport"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *UserRoleService) Assign(ctx context.Context, userID, roleID string, assignedBy *string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.assign")
	defer span.End()

	err := s.userRoleRepo.Assign(ctx, userID, roleID, assignedBy)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка назначения роли пользователю", err)
		return err
	}

	return nil
}
