package user_role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (s *APISuite) TestGetUserRolesSuccess() {
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
	expectedRoles := []*model.EnrichedRole{
		{Role: *role1, Permissions: []*model.Permission{}},
		{Role: *role2, Permissions: []*model.Permission{}},
	}

	req := &userRoleV1.GetUserRolesRequest{
		UserId: userID,
	}

	s.userRoleService.On("GetUserRoles", mock.Anything, mock.Anything).Return(expectedRoles, nil).Once()

	resp, err := s.api.GetUserRoles(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, len(expectedRoles))

	for i, role := range resp.Data {
		assert.Equal(s.T(), expectedRoles[i].Role.ID.String(), role.Role.Id)
		assert.Equal(s.T(), expectedRoles[i].Role.Name, role.Role.Name)
	}

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetUserRolesEmptyResult() {
	userID := "user123"
	expectedRoles := []*model.EnrichedRole{}

	req := &userRoleV1.GetUserRolesRequest{
		UserId: userID,
	}

	s.userRoleService.On("GetUserRoles", mock.Anything, mock.Anything).Return(expectedRoles, nil).Once()

	resp, err := s.api.GetUserRoles(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 0)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetUserRolesUserNotFound() {
	userID := "user123"

	req := &userRoleV1.GetUserRolesRequest{
		UserId: userID,
	}

	s.userRoleService.On("GetUserRoles", mock.Anything, mock.Anything).Return(nil, model.ErrUserRoleNotFound).Once()

	resp, err := s.api.GetUserRoles(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetUserRolesInternalError() {
	userID := "user123"

	req := &userRoleV1.GetUserRolesRequest{
		UserId: userID,
	}

	s.userRoleService.On("GetUserRoles", mock.Anything, mock.Anything).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.GetUserRoles(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}
