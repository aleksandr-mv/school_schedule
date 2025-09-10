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

func (s *APISuite) TestListPermissionsSuccess() {
	permission1 := &model.Permission{
		ID:       uuid.New(),
		Resource: "users",
		Action:   "read",
	}
	permission2 := &model.Permission{
		ID:       uuid.New(),
		Resource: "users",
		Action:   "write",
	}
	expectedPermissions := []*model.Permission{permission1, permission2}

	req := &permissionV1.ListPermissionsRequest{}

	s.permissionService.On("ListPermissions", mock.Anything, mock.Anything).Return(expectedPermissions, nil).Once()

	resp, err := s.api.ListPermissions(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 2)

	for i, permission := range resp.Data {
		assert.Equal(s.T(), expectedPermissions[i].ID.String(), permission.Id)
		assert.Equal(s.T(), expectedPermissions[i].Resource, permission.Resource)
		assert.Equal(s.T(), expectedPermissions[i].Action, permission.Action)
	}

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestListPermissionsWithResourceFilter() {
	permission1 := &model.Permission{
		ID:       uuid.New(),
		Resource: "users",
		Action:   "read",
	}
	expectedPermissions := []*model.Permission{permission1}
	resourceFilter := "users"

	req := &permissionV1.ListPermissionsRequest{
		ResourceFilter: &resourceFilter,
	}

	s.permissionService.On("ListPermissions", mock.Anything, mock.Anything).Return(expectedPermissions, nil).Once()

	resp, err := s.api.ListPermissions(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 1)

	assert.Equal(s.T(), expectedPermissions[0].ID.String(), resp.Data[0].Id)
	assert.Equal(s.T(), expectedPermissions[0].Resource, resp.Data[0].Resource)
	assert.Equal(s.T(), expectedPermissions[0].Action, resp.Data[0].Action)

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestListPermissionsWithActionFilter() {
	permission1 := &model.Permission{
		ID:       uuid.New(),
		Resource: "users",
		Action:   "read",
	}
	expectedPermissions := []*model.Permission{permission1}
	actionFilter := "read"

	req := &permissionV1.ListPermissionsRequest{
		ActionFilter: &actionFilter,
	}

	s.permissionService.On("ListPermissions", mock.Anything, mock.Anything).Return(expectedPermissions, nil).Once()

	resp, err := s.api.ListPermissions(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 1)

	assert.Equal(s.T(), expectedPermissions[0].ID.String(), resp.Data[0].Id)
	assert.Equal(s.T(), expectedPermissions[0].Resource, resp.Data[0].Resource)
	assert.Equal(s.T(), expectedPermissions[0].Action, resp.Data[0].Action)

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestListPermissionsEmptyResult() {
	expectedPermissions := []*model.Permission{}
	resourceFilter := "nonexistent"

	req := &permissionV1.ListPermissionsRequest{
		ResourceFilter: &resourceFilter,
	}

	s.permissionService.On("ListPermissions", mock.Anything, mock.Anything).Return(expectedPermissions, nil).Once()

	resp, err := s.api.ListPermissions(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 0)

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestListPermissionsInternalError() {
	req := &permissionV1.ListPermissionsRequest{}

	s.permissionService.On("ListPermissions", mock.Anything, mock.Anything).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.ListPermissions(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}
