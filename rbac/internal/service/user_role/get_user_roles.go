package user_role

import (
	"context"
	"sync"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/errreport"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/tracing"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *UserRoleService) GetUserRoles(ctx context.Context, userID string) ([]*model.EnrichedRole, error) {
	ctx, span := tracing.StartSpan(ctx, "rbac.service.get_user_roles")
	defer span.End()

	roleIDs, err := s.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения ролей пользователя из репозитория", err)
		return nil, err
	}

	enrichedRoles := make([]*model.EnrichedRole, len(roleIDs))
	errors := make([]error, len(roleIDs))

	const maxConcurrency = 10
	semaphore := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	for i, roleID := range roleIDs {
		wg.Add(1)
		go func(i int, roleID string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			enrichedRoles[i], errors[i] = s.roleService.Get(ctx, roleID)
		}(i, roleID)
	}

	wg.Wait()

	for _, err := range errors {
		if err != nil {
			errreport.Report(ctx, "❌ [Service] Ошибка получения обогащенной роли", err)
			return nil, err
		}
	}

	return enrichedRoles, nil
}
