package user

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

func (s *UserService) Register(ctx context.Context, createUser *model.CreateUser) (*model.User, error) {
	if err := createUser.Validate(); err != nil {
		logger.Error(ctx,
			"❌ [Service] Невалидные данные для регистрации",
			zap.Error(err),
			zap.String("operation", "user.Service.Register"),
		)

		return nil, model.ErrBadRequest
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(createUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx,
			"❌ [Service] Ошибка хэширования пароля",
			zap.Error(err),
			zap.String("operation", "user.Service.Register"),
		)

		return nil, model.ErrInternal
	}

	user := model.User{
		Login:        createUser.Login,
		Email:        createUser.Email,
		PasswordHash: string(hashedBytes),
	}

	createdUser, err := s.userRepository.Create(ctx, user)
	if err != nil {
		logger.Error(ctx,
			"❌ [Service] Ошибка создания пользователя в БД",
			zap.Error(err),
			zap.String("operation", "user.Service.Register"),
		)

		return nil, err
	}

	return createdUser, nil
}
