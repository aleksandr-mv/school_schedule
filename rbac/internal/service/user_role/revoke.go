package user_role

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
)

func (s *UserRoleService) Revoke(ctx context.Context, userID, roleID string) error {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.revoke_role")
	defer span.End()

	err := s.userRoleRepo.Revoke(ctx, userID, roleID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка отзыва роли у пользователя", err)
		return err
	}

	return nil
}
