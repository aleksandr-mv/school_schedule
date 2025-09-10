package permission_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestCheckPermissionSuccess() {
	userID := "user123"
	resource := "users"
	action := "read"

	permissionID := uuid.New()
	permission := &model.Permission{
		ID:       permissionID,
		Resource: resource,
		Action:   action,
	}

	roleID := uuid.New()
	roles := []*model.Role{
		{
			ID:        roleID,
			Name:      "admin",
			CreatedAt: time.Now(),
		},
	}

	s.permissionRepository.On("GetByResourceAndAction", mock.Anything, resource, action).Return(permission, nil)
	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return(roles, nil)
	s.rolePermissionRepository.On("HasPermission", mock.Anything, roleID.String(), permissionID.String()).Return(true, nil)

	err := s.service.CheckPermission(s.ctx, userID, resource, action)

	assert.NoError(s.T(), err)

	s.permissionRepository.AssertExpectations(s.T())
	s.userRoleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCheckPermissionNotFound() {
	userID := "user123"
	resource := "users"
	action := "read"

	s.permissionRepository.On("GetByResourceAndAction", mock.Anything, resource, action).Return(nil, model.ErrPermissionNotFound)

	err := s.service.CheckPermission(s.ctx, userID, resource, action)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionDenied, err)

	s.permissionRepository.AssertExpectations(s.T())
	s.userRoleRepository.AssertNotCalled(s.T(), "GetUserRoles")
	s.rolePermissionRepository.AssertNotCalled(s.T(), "HasPermission")
}

func (s *ServiceSuite) TestCheckPermissionUserHasNoRoles() {
	userID := "user123"
	resource := "users"
	action := "read"

	permissionID := uuid.New()
	permission := &model.Permission{
		ID:       permissionID,
		Resource: resource,
		Action:   action,
	}

	s.permissionRepository.On("GetByResourceAndAction", mock.Anything, resource, action).Return(permission, nil)
	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return([]*model.Role{}, nil)

	err := s.service.CheckPermission(s.ctx, userID, resource, action)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionDenied, err)

	s.permissionRepository.AssertExpectations(s.T())
	s.userRoleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertNotCalled(s.T(), "HasPermission")
}

func (s *ServiceSuite) TestCheckPermissionUserRolesNotFound() {
	userID := "user123"
	resource := "users"
	action := "read"

	permissionID := uuid.New()
	permission := &model.Permission{
		ID:       permissionID,
		Resource: resource,
		Action:   action,
	}

	s.permissionRepository.On("GetByResourceAndAction", mock.Anything, resource, action).Return(permission, nil)
	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return(nil, model.ErrUserRoleNotFound)

	err := s.service.CheckPermission(s.ctx, userID, resource, action)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrUserRoleNotFound, err)

	s.permissionRepository.AssertExpectations(s.T())
	s.userRoleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertNotCalled(s.T(), "HasPermission")
}

func (s *ServiceSuite) TestCheckPermissionNoPermission() {
	userID := "user123"
	resource := "users"
	action := "read"

	permissionID := uuid.New()
	permission := &model.Permission{
		ID:       permissionID,
		Resource: resource,
		Action:   action,
	}

	roleID := uuid.New()
	roles := []*model.Role{
		{
			ID:        roleID,
			Name:      "user",
			CreatedAt: time.Now(),
		},
	}

	s.permissionRepository.On("GetByResourceAndAction", mock.Anything, resource, action).Return(permission, nil)
	s.userRoleRepository.On("GetUserRoles", mock.Anything, userID).Return(roles, nil)
	s.rolePermissionRepository.On("HasPermission", mock.Anything, roleID.String(), permissionID.String()).Return(false, nil)

	err := s.service.CheckPermission(s.ctx, userID, resource, action)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrPermissionDenied, err)

	s.permissionRepository.AssertExpectations(s.T())
	s.userRoleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertExpectations(s.T())
}
