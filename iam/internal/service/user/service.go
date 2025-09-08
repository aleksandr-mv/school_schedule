package user

import (
	"github.com/aleksandr-mv/school_schedule/iam/internal/repository"
	def "github.com/aleksandr-mv/school_schedule/iam/internal/service"
)

var _ def.UserServiceInterface = (*UserService)(nil)

type UserService struct {
	userRepository         repository.UserRepository
	notificationRepository repository.NotificationRepository
}

func NewService(userRepository repository.UserRepository, notificationRepository repository.NotificationRepository) *UserService {
	return &UserService{
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
	}
}
