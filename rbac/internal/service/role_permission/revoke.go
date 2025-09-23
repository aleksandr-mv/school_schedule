package role_permission

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
)

func (s *RolePermissionService) Revoke(ctx context.Context, roleID, permissionID string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.revoke_permission_from_role")
	defer span.End()

	err := s.rolePermissionRepo.Revoke(ctx, roleID, permissionID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка отзыва права у роли", err)
		return err
	}

	return nil
}
