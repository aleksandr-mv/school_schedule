package user

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, id.String())
	if err != nil {
		logger.Error(ctx,
			"❌ [Service] Ошибка получения пользователя",
			zap.Error(err),
			zap.String("operation", "user.Service.GetUser"),
			zap.String("user_id", id.String()),
		)

		return nil, err
	}

	notificationMethods, err := s.notificationRepository.GetByUser(ctx, user.ID)
	if err != nil {
		logger.Error(ctx,
			"❌ [Service] Ошибка получения методов уведомлений",
			zap.Error(err),
			zap.String("operation", "user.Service.GetUser"),
			zap.String("user_id", id.String()),
		)
		return nil, err
	}

	user.NotificationMethods = notificationMethods

	return user, nil
}
