package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/aleksandr-mv/school_schedule/iam/internal/repository/mocks"
	"github.com/aleksandr-mv/school_schedule/iam/internal/service/user"
	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context // nolint:containedctx

	userRepository         *mocks.UserRepository
	notificationRepository *mocks.NotificationRepository

	service *user.UserService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	if err := logger.InitDefault(); err != nil {
		panic(err)
	}

	s.userRepository = mocks.NewUserRepository(s.T())
	s.notificationRepository = mocks.NewNotificationRepository(s.T())

	s.service = user.NewService(s.userRepository, s.notificationRepository)
}

func (s *ServiceSuite) SetupTest() {
	s.userRepository.ExpectedCalls = nil
	s.notificationRepository.ExpectedCalls = nil
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
