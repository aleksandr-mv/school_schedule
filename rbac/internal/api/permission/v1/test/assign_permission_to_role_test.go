package permission_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (s *APISuite) TestAssignPermissionToRoleSuccess() {
	roleID := uuid.New()
	permissionID := uuid.New()

	req := &permissionV1.AssignPermissionToRoleRequest{
		RoleId:       roleID.String(),
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("AssignPermissionToRole", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	resp, err := s.api.AssignPermissionToRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignPermissionToRolePermissionAlreadyAssigned() {
	roleID := uuid.New()
	permissionID := uuid.New()

	req := &permissionV1.AssignPermissionToRoleRequest{
		RoleId:       roleID.String(),
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("AssignPermissionToRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrPermissionAlreadyAssigned).Once()

	resp, err := s.api.AssignPermissionToRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.AlreadyExists, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignPermissionToRoleRoleNotFound() {
	roleID := uuid.New()
	permissionID := uuid.New()

	req := &permissionV1.AssignPermissionToRoleRequest{
		RoleId:       roleID.String(),
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("AssignPermissionToRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.AssignPermissionToRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignPermissionToRolePermissionNotFound() {
	roleID := uuid.New()
	permissionID := uuid.New()

	req := &permissionV1.AssignPermissionToRoleRequest{
		RoleId:       roleID.String(),
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("AssignPermissionToRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrPermissionNotFound).Once()

	resp, err := s.api.AssignPermissionToRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignPermissionToRoleInternalError() {
	roleID := uuid.New()
	permissionID := uuid.New()

	req := &permissionV1.AssignPermissionToRoleRequest{
		RoleId:       roleID.String(),
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("AssignPermissionToRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrInternal).Once()

	resp, err := s.api.AssignPermissionToRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}
