package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestRevokeRoleSuccess() {
	userID := "user123"
	roleID := "role456"

	s.userRoleRepository.On("RevokeRole", mock.Anything, userID, roleID).Return(nil)

	err := s.service.RevokeRole(s.ctx, userID, roleID)

	assert.NoError(s.T(), err)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokeRoleNotAssigned() {
	userID := "user123"
	roleID := "role456"

	s.userRoleRepository.On("RevokeRole", mock.Anything, userID, roleID).Return(model.ErrRoleNotAssigned)

	err := s.service.RevokeRole(s.ctx, userID, roleID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrRoleNotAssigned, err)

	s.userRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRevokeRoleRepositoryError() {
	userID := "user123"
	roleID := "role456"

	s.userRoleRepository.On("RevokeRole", mock.Anything, userID, roleID).Return(model.ErrInternal)

	err := s.service.RevokeRole(s.ctx, userID, roleID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.userRoleRepository.AssertExpectations(s.T())
}
