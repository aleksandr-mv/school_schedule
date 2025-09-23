package user_producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	def "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/errreport"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/kafka"
	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

var _ def.UserProducerService = (*service)(nil)

type service struct {
	producer kafka.Producer
}

func NewService(producer kafka.Producer) def.UserProducerService {
	return &service{
		producer: producer,
	}
}

func (s *service) ProduceUserCreated(ctx context.Context, event model.UserCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		errreport.Report(ctx, "❌ [Producer] Ошибка кодирования UserCreated", err)
		return fmt.Errorf("encode user created: %w", err)
	}

	if err = s.producer.Send(ctx, []byte(event.UserID.String()), payload); err != nil {
		errreport.Report(ctx, "❌ [Producer] Ошибка отправки UserCreated", err)
		return fmt.Errorf("send user created to kafka: %w", err)
	}

	logger.Info(ctx, "📤 Отправлено событие UserCreated")

	return nil
}
