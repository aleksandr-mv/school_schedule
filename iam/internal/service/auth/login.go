package auth

import (
	"context"
	"time"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(ctx context.Context, credentials *model.LoginCredentials) (uuid.UUID, error) {
	if err := credentials.Validate(); err != nil {
		errreport.Report(ctx, "❌ [Service] Невалидные учетные данные", err)
		return uuid.Nil, model.ErrInvalidCredentials
	}

	user, err := s.userRepository.Get(ctx, credentials.Login)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения пользователя", err)
		return uuid.Nil, model.ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		logger.Error(ctx,
			"⚠️ [Service] Неверный пароль",
			zap.String("operation", "auth.Service.Login"),
			zap.String("user_id", user.ID.String()),
		)

		return uuid.Nil, model.ErrInvalidCredentials
	}

	notificationMethods, err := s.notificationRepository.GetByUser(ctx, user.ID)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка получения методов уведомлений", err)
		return uuid.Nil, err
	}

	user.NotificationMethods = notificationMethods

	expiresAt := time.Now().Add(s.sessionTTL)

	sessionID, err := s.sessionRepository.Create(ctx, *user, expiresAt)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка создания сессии", err)
		return uuid.Nil, model.ErrFailedToCreateSession
	}

	return sessionID, nil
}
