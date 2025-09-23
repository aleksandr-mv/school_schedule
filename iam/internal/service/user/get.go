package user

import (
	"context"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/google/uuid"
)

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, id.String())
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения пользователя", err)
		return nil, err
	}

	notificationMethods, err := s.notificationRepository.GetByUser(ctx, user.ID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения методов уведомлений", err)
		return nil, err
	}

	user.NotificationMethods = notificationMethods

	return user, nil
}
