package user_role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestGetUserRolesSuccess() {
	userID := "user123"
	role1 := &model.Role{
		ID:        uuid.New(),
		Name:      "admin",
		CreatedAt: time.Now(),
	}
	role2 := &model.Role{
		ID:        uuid.New(),
		Name:      "user",
		CreatedAt: time.Now(),
	}
	expectedRoles := []*model.Role{role1, role2}

	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return(expectedRoles, nil)

	result, err := s.service.GetUserRoles(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), expectedRoles, result)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUserRolesEmptyResult() {
	userID := "user123"

	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return([]*model.Role{}, nil)

	result, err := s.service.GetUserRoles(s.ctx, userID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUserRolesNotFound() {
	userID := "user123"

	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return(nil, model.ErrUserRoleNotFound)

	result, err := s.service.GetUserRoles(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrUserRoleNotFound, err)
	assert.Nil(s.T(), result)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetUserRolesRepositoryError() {
	userID := "user123"

	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return(nil, model.ErrInternal)

	result, err := s.service.GetUserRoles(s.ctx, userID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)
	assert.Nil(s.T(), result)

	s.userRoleRepository.AssertExpectations(s.T())
}
