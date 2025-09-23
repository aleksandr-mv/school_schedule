package auth

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/google/uuid"
)

func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.sessionRepository.Delete(ctx, sessionID); err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка удаления сессии", err)
		return model.ErrFailedToDeleteSession
	}

	return nil
}
