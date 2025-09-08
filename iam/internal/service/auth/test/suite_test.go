package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/iam/internal/repository/mocks"
	"github.com/aleksandr-mv/school_schedule/iam/internal/service/auth"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	userRepository         *mocks.UserRepository
	notificationRepository *mocks.NotificationRepository
	sessionRepository      *mocks.SessionRepository

	service *auth.AuthService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.userRepository = mocks.NewUserRepository(s.T())
	s.notificationRepository = mocks.NewNotificationRepository(s.T())
	s.sessionRepository = mocks.NewSessionRepository(s.T())

	s.service = auth.NewService(s.userRepository, s.notificationRepository, s.sessionRepository, 24*time.Hour)
}

func (s *ServiceSuite) SetupTest() {
	s.userRepository.ExpectedCalls = nil
	s.notificationRepository.ExpectedCalls = nil
	s.sessionRepository.ExpectedCalls = nil
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
