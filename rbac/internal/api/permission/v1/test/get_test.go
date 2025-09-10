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

func (s *APISuite) TestGetPermissionSuccess() {
	permissionID := uuid.New()
	expectedPermission := &model.Permission{
		ID:       permissionID,
		Resource: "users",
		Action:   "read",
	}

	req := &permissionV1.GetPermissionRequest{
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("GetPermission", mock.Anything, mock.Anything).Return(expectedPermission, nil).Once()

	resp, err := s.api.GetPermission(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotNil(s.T(), resp.Permission)
	assert.Equal(s.T(), expectedPermission.ID.String(), resp.Permission.Id)
	assert.Equal(s.T(), expectedPermission.Resource, resp.Permission.Resource)
	assert.Equal(s.T(), expectedPermission.Action, resp.Permission.Action)

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetPermissionNotFound() {
	permissionID := uuid.New()

	req := &permissionV1.GetPermissionRequest{
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("GetPermission", mock.Anything, mock.Anything).Return(nil, model.ErrPermissionNotFound).Once()

	resp, err := s.api.GetPermission(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetPermissionInternalError() {
	permissionID := uuid.New()

	req := &permissionV1.GetPermissionRequest{
		PermissionId: permissionID.String(),
	}

	s.permissionService.On("GetPermission", mock.Anything, mock.Anything).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.GetPermission(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}
