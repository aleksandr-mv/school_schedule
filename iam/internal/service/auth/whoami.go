package auth

import (
	"context"
	"time"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/errreport"
	"github.com/google/uuid"
)

func (s *AuthService) Whoami(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error) {
	iam, err := s.sessionRepository.Get(ctx, sessionID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения сессии", err)
		return nil, err
	}

	if iam.Session.ExpiresAt.Before(time.Now()) {
		return nil, model.ErrSessionExpired
	}

	roles, err := s.rbacClient.GetUserRoles(ctx, iam.User.ID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения ролей пользователя", err)
		roles = []*model.RoleWithPermissions{}
	}

	iam.RolesWithPermissions = roles

	return iam, nil
}
