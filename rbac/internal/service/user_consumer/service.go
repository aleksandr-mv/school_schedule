package user_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	def "github.com/aleksandr-mv/school_schedule/rbac/internal/service"
)

var _ def.UserConsumerService = (*service)(nil)

type service struct {
	userCreatedConsumer kafka.Consumer
	userRoleService     def.UserRoleServiceInterface
}

func NewService(
	userCreatedConsumer kafka.Consumer,
	userRoleService def.UserRoleServiceInterface,
) *service {
	return &service{
		userCreatedConsumer: userCreatedConsumer,
		userRoleService:     userRoleService,
	}
}

func (s *service) Run(ctx context.Context) error {
	logger.Info(ctx, "🚀 Запуск UserCreated Consumer")

	if err := s.userCreatedConsumer.Consume(ctx, s.UserCreatedHandler); err != nil {
		logger.Error(ctx, "❌ Ошибка в UserCreated consumer", zap.Error(err))
		return err
	}

	return nil
}
