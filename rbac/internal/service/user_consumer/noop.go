package user_consumer

import (
	"context"

	def "github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/service"
)

type noOpService struct{}

// NewNoOpService создает no-op реализацию UserConsumerService
// Используется когда Kafka отключен
func NewNoOpService() def.UserConsumerService {
	return &noOpService{}
}

func (n *noOpService) Run(ctx context.Context) error {
	// No-op: ничего не делаем, когда Kafka отключен
	return nil
}
