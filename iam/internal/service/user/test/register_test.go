package user_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
)

func (s *ServiceSuite) TestRegisterSuccess() {
	userID := uuid.New()

	login := "testuser123"
	email := "test@example.com"
	password := "password123456"

	expectedUser := &model.User{
		ID:           userID,
		Login:        "testuser123",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
	}

	s.userRepository.On("Create", mock.Anything, mock.AnythingOfType("model.User")).Return(expectedUser, nil)
	s.userProducerService.On("ProduceUserCreated", mock.Anything, mock.AnythingOfType("model.UserCreated")).Return(nil)

	result, err := s.service.Register(s.ctx, login, email, password)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.ID, result.ID)
	assert.Equal(s.T(), expectedUser.Login, result.Login)
	assert.Equal(s.T(), expectedUser.Email, result.Email)

	s.userRepository.AssertExpectations(s.T())
	s.userProducerService.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegisterUserAlreadyExists() {
	login := "existinguser"
	email := "existing@example.com"
	password := "password123456"

	s.userRepository.On("Create", mock.Anything, mock.AnythingOfType("model.User")).Return(nil, model.ErrUserAlreadyExists)

	result, err := s.service.Register(s.ctx, login, email, password)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrUserAlreadyExists, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertExpectations(s.T())
	// Producer не должен быть вызван, если пользователь не создан
	s.userProducerService.AssertNotCalled(s.T(), "ProduceUserCreated")
}

func (s *ServiceSuite) TestRegisterRepositoryError() {
	login := "testuser123"
	email := "test@example.com"
	password := "password123456"

	s.userRepository.On("Create", mock.Anything, mock.AnythingOfType("model.User")).Return(nil, model.ErrInternal)

	result, err := s.service.Register(s.ctx, login, email, password)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertExpectations(s.T())
}
