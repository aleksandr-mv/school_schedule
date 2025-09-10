package permission_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestAssignPermissionToRoleSuccess() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("AssignPermissionToRole", mock.Anything, roleID, permissionID).Return(nil)

	err := s.service.AssignPermissionToRole(s.ctx, roleID, permissionID)

	assert.NoError(s.T(), err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAssignPermissionToRoleAlreadyAssigned() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("AssignPermissionToRole", mock.Anything, roleID, permissionID).Return(model.ErrPermissionAlreadyAssigned)

	err := s.service.AssignPermissionToRole(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionAlreadyAssigned, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestAssignPermissionToRoleRepositoryError() {
	roleID := "role123"
	permissionID := "permission456"

	s.rolePermissionRepository.On("AssignPermissionToRole", mock.Anything, roleID, permissionID).Return(model.ErrInternal)

	err := s.service.AssignPermissionToRole(s.ctx, roleID, permissionID)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.rolePermissionRepository.AssertExpectations(s.T())
}
