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

	createUser := &model.CreateUser{
		Login:    "testuser123",
		Email:    "test@example.com",
		Password: "password123456",
	}

	expectedUser := &model.User{
		ID:           userID,
		Login:        "testuser123",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    nil,
	}

	s.userRepository.On("Create", mock.Anything, mock.AnythingOfType("model.User")).Return(expectedUser, nil)

	result, err := s.service.Register(s.ctx, createUser)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedUser.ID, result.ID)
	assert.Equal(s.T(), expectedUser.Login, result.Login)
	assert.Equal(s.T(), expectedUser.Email, result.Email)

	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegisterUserAlreadyExists() {
	createUser := &model.CreateUser{
		Login:    "existinguser",
		Email:    "existing@example.com",
		Password: "password123456",
	}

	s.userRepository.On("Create", mock.Anything, mock.AnythingOfType("model.User")).Return(nil, model.ErrUserAlreadyExists)

	result, err := s.service.Register(s.ctx, createUser)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrUserAlreadyExists, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegisterInvalidData() {
	createUser := &model.CreateUser{
		Login:    "", // Невалидный логин
		Email:    "test@example.com",
		Password: "password123456",
	}

	result, err := s.service.Register(s.ctx, createUser)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrBadRequest, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestRegisterRepositoryError() {
	createUser := &model.CreateUser{
		Login:    "testuser123",
		Email:    "test@example.com",
		Password: "password123456",
	}

	s.userRepository.On("Create", mock.Anything, mock.AnythingOfType("model.User")).Return(nil, model.ErrInternal)

	result, err := s.service.Register(s.ctx, createUser)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.userRepository.AssertExpectations(s.T())
}
