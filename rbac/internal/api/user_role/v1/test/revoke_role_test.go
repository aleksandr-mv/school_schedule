package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (s *APISuite) TestRevokeRoleSuccess() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRoleRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("RevokeRole", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	resp, err := s.api.RevokeRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeRoleRoleNotAssigned() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRoleRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("RevokeRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotAssigned).Once()

	resp, err := s.api.RevokeRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.FailedPrecondition, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeRoleUserNotFound() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRoleRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("RevokeRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.RevokeRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeRoleRoleNotFound() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRoleRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("RevokeRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.RevokeRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeRoleInternalError() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRoleRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("RevokeRole", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrInternal).Once()

	resp, err := s.api.RevokeRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}
