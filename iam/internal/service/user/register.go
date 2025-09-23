package user

import (
	"context"
	"fmt"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/errreport"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) Register(ctx context.Context, login, email, password string) (*model.User, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка хэширования пароля", err)
		return nil, model.ErrInternal
	}

	user := model.User{
		Login:        login,
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	createdUser, err := s.userRepository.Create(ctx, user)
	if err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка создания пользователя в БД", err)
		return nil, err
	}

	defaultRoleID := model.DefaultRoleID
	if err := s.userProducerService.ProduceUserCreated(ctx, model.NewUserCreated(createdUser, defaultRoleID)); err != nil {
		errreport.Report(ctx, "❌ [Service] Ошибка отправки события UserCreated", err)

		if deleteErr := s.userRepository.Delete(ctx, createdUser.ID); deleteErr != nil {
			errreport.Report(ctx, "❌ [Service] Критическая ошибка: не удалось удалить пользователя при откате", deleteErr)
		}

		return nil, fmt.Errorf("failed to send user created event: %w", err)
	}

	return createdUser, nil
}
