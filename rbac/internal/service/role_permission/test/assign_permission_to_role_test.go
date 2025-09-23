package role_permission_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestAssignSuccess() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("Assign", mock.Anything, roleID, permissionID).Return(nil)

	err := s.service.Assign(s.ctx, roleID, permissionID)

	assert.NoError(s.T(), err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAssignAlreadyAssigned() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("Assign", mock.Anything, roleID, permissionID).Return(model.ErrPermissionAlreadyAssigned)

	err := s.service.Assign(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionAlreadyAssigned, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAssignRepositoryError() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("Assign", mock.Anything, roleID, permissionID).Return(model.ErrInternal)

	err := s.service.Assign(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}
