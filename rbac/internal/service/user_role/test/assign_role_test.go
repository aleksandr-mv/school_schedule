package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestAssignRoleSuccess() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	s.userRoleRepository.On("AssignRole", mock.Anything, userID, roleID, &assignedBy).Return(nil)

	err := s.service.AssignRole(s.ctx, userID, roleID, &assignedBy)

	assert.NoError(s.T(), err)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAssignRoleSuccessWithoutAssignedBy() {
	userID := "user123"
	roleID := "role456"

	s.userRoleRepository.On("AssignRole", mock.Anything, userID, roleID, (*string)(nil)).Return(nil)

	err := s.service.AssignRole(s.ctx, userID, roleID, nil)

	assert.NoError(s.T(), err)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAssignRoleAlreadyAssigned() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	s.userRoleRepository.On("AssignRole", mock.Anything, userID, roleID, &assignedBy).Return(model.ErrRoleAlreadyAssigned)

	err := s.service.AssignRole(s.ctx, userID, roleID, &assignedBy)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleAlreadyAssigned, err)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAssignRoleRepositoryError() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	s.userRoleRepository.On("AssignRole", mock.Anything, userID, roleID, &assignedBy).Return(model.ErrInternal)

	err := s.service.AssignRole(s.ctx, userID, roleID, &assignedBy)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.userRoleRepository.AssertExpectations(s.T())
}
