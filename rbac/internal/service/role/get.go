package role

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/tracing"
	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *RoleService) Get(ctx context.Context, id string) (*model.EnrichedRole, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_role")
	defer span.End()

	enrichedRole, err := s.enrichedRoleRepo.Get(ctx, id)
	if err == nil {
		return enrichedRole, nil
	}

	role, err := s.roleRepo.Get(ctx, id)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения роли из репозитория", err)
		return nil, err
	}

	permissions, err := s.rolePermissionRepo.GetRolePermissions(ctx, id)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения прав роли", err)
		return nil, err
	}

	enrichedRole = &model.EnrichedRole{
		Role:        *role,
		Permissions: permissions,
	}

	expiresAt := time.Now().Add(s.enrichedRoleTTL)
	if err := s.enrichedRoleRepo.Set(ctx, enrichedRole, expiresAt); err != nil {
		logger.Warn(ctx, "⚠️ [Service] Не удалось кэшировать роль", zap.Error(err))
	}

	return enrichedRole, nil
}
