package user_consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/kafka/consumer"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *service) UserCreatedHandler(ctx context.Context, msg consumer.Message) error {
	var event model.UserCreated
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.Error(ctx, "❌ Ошибка декодирования UserCreated", zap.Error(err))
		return fmt.Errorf("decode user created: %w", err)
	}

	logger.Info(ctx, "📥 Получено событие UserCreated",
		zap.String("topic", msg.Topic))

	if err := s.userRoleService.Assign(ctx, event.UserID.String(), event.RoleID, nil); err != nil {
		logger.Error(ctx, "❌ Ошибка назначения роли", zap.Error(err))
		return fmt.Errorf("assign role: %w", err)
	}

	return nil
}
