package auth

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.sessionRepository.Delete(ctx, sessionID); err != nil {
		logger.Error(ctx,
			"❌ [Service] Ошибка удаления сессии",
			zap.Error(err),
			zap.String("operation", "auth.Service.Logout"),
			zap.String("session_id", sessionID.String()),
		)

		return model.ErrFailedToDeleteSession
	}

	return nil
}
