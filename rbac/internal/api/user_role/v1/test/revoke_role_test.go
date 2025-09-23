package user_role_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	userRoleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user_role/v1"
)

func (s *APISuite) TestRevokeSuccess() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("Revoke", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	resp, err := s.api.Revoke(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeRoleNotAssigned() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("Revoke", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotAssigned).Once()

	resp, err := s.api.Revoke(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.FailedPrecondition, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeUserNotFound() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("Revoke", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.Revoke(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeRoleNotFound() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("Revoke", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.Revoke(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeInternalError() {
	userID := "user123"
	roleID := "role456"

	req := &userRoleV1.RevokeRequest{
		UserId: userID,
		RoleId: roleID,
	}

	s.userRoleService.On("Revoke", mock.Anything, mock.Anything, mock.Anything).Return(model.ErrInternal).Once()

	resp, err := s.api.Revoke(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestRevokeValidation_InvalidUserID() {
	req := &userRoleV1.RevokeRequest{
		UserId: "invalid-uuid",
		RoleId: uuid.New().String(),
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value must be a valid UUID")
}

func (s *APISuite) TestRevokeValidation_InvalidRoleID() {
	req := &userRoleV1.RevokeRequest{
		UserId: uuid.New().String(),
		RoleId: "invalid-uuid",
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value must be a valid UUID")
}

func (s *APISuite) TestRevokeValidation_ValidRequest() {
	req := &userRoleV1.RevokeRequest{
		UserId: uuid.New().String(),
		RoleId: uuid.New().String(),
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}
