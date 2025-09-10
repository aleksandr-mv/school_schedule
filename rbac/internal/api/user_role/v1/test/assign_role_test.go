package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (s *APISuite) TestAssignRoleSuccess() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("AssignRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	resp, err := s.api.AssignRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignRoleSuccessWithoutAssignedBy() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: nil,
	}

	s.userRoleService.On("AssignRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	resp, err := s.api.AssignRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignRoleRoleAlreadyAssigned() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("AssignRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleAlreadyAssigned).Once()

	resp, err := s.api.AssignRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.AlreadyExists, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignRoleUserNotFound() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("AssignRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.AssignRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignRoleRoleNotFound() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("AssignRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.AssignRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestAssignRoleInternalError() {
	userID := "user123"
	roleID := "role456"
	assignedBy := "admin123"

	req := &userRoleV1.AssignRoleRequest{
		UserId:     userID,
		RoleId:     roleID,
		AssignedBy: &assignedBy,
	}

	s.userRoleService.On("AssignRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.ErrInternal).Once()

	resp, err := s.api.AssignRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}
