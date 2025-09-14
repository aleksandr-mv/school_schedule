package user_producer

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	def "github.com/aleksandr-mv/school_schedule/iam/internal/service"
)

type noOpService struct{}

// NewNoOpService создает no-op реализацию UserProducerService
// Используется когда Kafka отключен
func NewNoOpService() def.UserProducerService {
	return &noOpService{}
}

func (n *noOpService) ProduceUserCreated(ctx context.Context, event model.UserCreated) error {
	return nil
}
