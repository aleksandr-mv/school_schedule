package auth

import (
	"time"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/client/grpc"
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/repository"
	def "github.com/Alexander-Mandzhiev/school_schedule/iam/internal/service"
)

var _ def.AuthService = (*AuthService)(nil)

type AuthService struct {
	userRepository         repository.UserRepository
	notificationRepository repository.NotificationRepository
	sessionRepository      repository.SessionRepository
	rbacClient             grpc.RBACClient
	sessionTTL             time.Duration
}

func NewService(
	userRepository repository.UserRepository,
	notificationRepository repository.NotificationRepository,
	sessionRepository repository.SessionRepository,
	rbacClient grpc.RBACClient,
	sessionTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
		sessionRepository:      sessionRepository,
		rbacClient:             rbacClient,
		sessionTTL:             sessionTTL,
	}
}
