package role_permission

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/errreport"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
)

func (s *RolePermissionService) Assign(ctx context.Context, roleID, permissionID string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.assign_permission_to_role")
	defer span.End()

	err := s.rolePermissionRepo.Assign(ctx, roleID, permissionID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка назначения права роли", err)
		return err
	}

	return nil
}
