package permission

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *PermissionService) List(ctx context.Context) ([]*model.Permission, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.list_permissions")
	defer span.End()

	permissions, err := s.permissionRepo.List(ctx)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения списка прав доступа из репозитория", err)
		return nil, err
	}

	return permissions, nil
}
