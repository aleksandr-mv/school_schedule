package whoami

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
)

func (s *WhoAMIService) Whoami(ctx context.Context, sessionID uuid.UUID) (*model.WhoAMI, error) {
	iam, err := s.sessionRepository.Get(ctx, sessionID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения сессии", err)
		return nil, err
	}
	if iam.Session.ExpiresAt.Before(time.Now()) {
		return nil, model.ErrSessionExpired
	}

	return iam, nil
}
