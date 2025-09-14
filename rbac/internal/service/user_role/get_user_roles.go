package user_role

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *UserRoleService) GetUserRoles(ctx context.Context, userID string) ([]*model.EnrichedRole, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_user_roles")
	defer span.End()

	roleIDs, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		logger.Error(ctx, "❌ [Service] Ошибка получения ролей пользователя из репозитория", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Параллельное получение обогащенных ролей для избежания N+1 проблемы
	enrichedRoles := make([]*model.EnrichedRole, len(roleIDs))
	errors := make([]error, len(roleIDs))

	// Семафор для ограничения количества параллельных запросов
	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	for i, roleID := range roleIDs {
		wg.Add(1)
		go func(i int, roleID string) {
			defer wg.Done()

			// Захватываем семафор
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			enrichedRoles[i], errors[i] = s.roleService.Get(ctx, roleID)
		}(i, roleID)
	}

	wg.Wait()

	// Проверяем ошибки
	for i, err := range errors {
		if err != nil {
			logger.Error(ctx, "❌ [Service] Ошибка получения обогащенной роли",
				zap.String("roleID", roleIDs[i]), zap.Error(err))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	return enrichedRoles, nil
}
