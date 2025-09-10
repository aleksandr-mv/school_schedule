package user_role_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	userRoleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
)

func (s *APISuite) TestGetRoleUsersSuccess() {
	roleID := "role123"
	limit := int32(10)
	cursor := "cursor123"
	nextCursor := "nextCursor456"

	expectedUsers := []string{"user1", "user2", "user3"}
	expectedTotal := int32(3)

	req := &userRoleV1.GetRoleUsersRequest{
		RoleId: roleID,
		Limit:  &limit,
		Cursor: &cursor,
	}

	s.userRoleService.On("GetRoleUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedUsers, expectedTotal, &nextCursor, nil).Once()

	resp, err := s.api.GetRoleUsers(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.UserIds, len(expectedUsers))
	assert.Equal(s.T(), expectedTotal, resp.TotalCount)
	assert.NotNil(s.T(), resp.NextCursor)
	assert.Equal(s.T(), nextCursor, *resp.NextCursor)

	for i, user := range resp.UserIds {
		assert.Equal(s.T(), expectedUsers[i], user)
	}

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleUsersSuccessWithoutCursor() {
	roleID := "role123"
	limit := int32(10)

	expectedUsers := []string{"user1", "user2", "user3"}
	expectedTotal := int32(3)

	req := &userRoleV1.GetRoleUsersRequest{
		RoleId: roleID,
		Limit:  &limit,
		Cursor: nil,
	}

	s.userRoleService.On("GetRoleUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedUsers, expectedTotal, nil, nil).Once()

	resp, err := s.api.GetRoleUsers(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.UserIds, len(expectedUsers))
	assert.Equal(s.T(), expectedTotal, resp.TotalCount)
	assert.Nil(s.T(), resp.NextCursor)

	for i, user := range resp.UserIds {
		assert.Equal(s.T(), expectedUsers[i], user)
	}

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleUsersEmptyResult() {
	roleID := "role123"
	limit := int32(10)

	expectedUsers := []string{}
	expectedTotal := int32(0)

	req := &userRoleV1.GetRoleUsersRequest{
		RoleId: roleID,
		Limit:  &limit,
		Cursor: nil,
	}

	s.userRoleService.On("GetRoleUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedUsers, expectedTotal, nil, nil).Once()

	resp, err := s.api.GetRoleUsers(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.UserIds, 0)
	assert.Equal(s.T(), expectedTotal, resp.TotalCount)
	assert.Nil(s.T(), resp.NextCursor)

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleUsersRoleNotFound() {
	roleID := "role123"
	limit := int32(10)

	req := &userRoleV1.GetRoleUsersRequest{
		RoleId: roleID,
		Limit:  &limit,
		Cursor: nil,
	}

	s.userRoleService.On("GetRoleUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, int32(0), nil, model.ErrRoleNotFound).Once()

	resp, err := s.api.GetRoleUsers(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleUsersInternalError() {
	roleID := "role123"
	limit := int32(10)

	req := &userRoleV1.GetRoleUsersRequest{
		RoleId: roleID,
		Limit:  &limit,
		Cursor: nil,
	}

	s.userRoleService.On("GetRoleUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, int32(0), nil, model.ErrInternal).Once()

	resp, err := s.api.GetRoleUsers(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.userRoleService.AssertExpectations(s.T())
}
