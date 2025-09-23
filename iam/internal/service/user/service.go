package user

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service"
)

var _ service.UserService = (*UserService)(nil)

type UserService struct {
	userRepository         repository.UserRepository
	notificationRepository repository.NotificationRepository
	userProducerService    service.UserProducerService
}

func NewService(
	userRepository repository.UserRepository,
	notificationRepository repository.NotificationRepository,
	userProducerService service.UserProducerService,
) *UserService {
	return &UserService{
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
		userProducerService:    userProducerService,
	}
}
