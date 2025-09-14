package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

func (s *AuthService) Whoami(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error) {
	iam, err := s.sessionRepository.Get(ctx, sessionID)
	if err != nil {
		logger.Error(ctx,
			"❌ [Service] Ошибка получения сессии",
			zap.Error(err),
			zap.String("operation", "auth.Service.Whoami"),
		)

		return nil, err
	}

	if iam.Session.ExpiresAt.Before(time.Now()) {
		return nil, model.ErrSessionExpired
	}

	roles, err := s.rbacClient.GetUserRoles(ctx, iam.User.ID)
	if err != nil {
		logger.Error(ctx,
			"❌ [Service] Ошибка получения ролей пользователя",
			zap.Error(err),
			zap.String("operation", "auth.Service.Whoami"),
			zap.String("userID", iam.User.ID.String()),
		)
		roles = []*model.RoleWithPermissions{}
	}

	iam.RolesWithPermissions = roles

	return iam, nil
}
