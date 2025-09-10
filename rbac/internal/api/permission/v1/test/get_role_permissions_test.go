package permission_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

func (s *APISuite) TestListPermissionsByRoleSuccess() {
	roleID := uuid.New()
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

	req := &permissionV1.ListPermissionsByRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: &commonV1.GetIdentifier_Id{
				Id: roleID.String(),
			},
		},
	}

	s.permissionService.On("ListPermissionsByRole", mock.Anything, mock.Anything).Return(expectedPermissions, nil).Once()

	resp, err := s.api.ListPermissionsByRole(s.ctx, req)

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

func (s *APISuite) TestListPermissionsByRoleNotFound() {
	roleID := uuid.New()

	req := &permissionV1.ListPermissionsByRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: &commonV1.GetIdentifier_Id{
				Id: roleID.String(),
			},
		},
	}

	s.permissionService.On("ListPermissionsByRole", mock.Anything, mock.Anything).Return(nil, model.ErrRoleNotFound).Once()

	resp, err := s.api.ListPermissionsByRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}

func (s *APISuite) TestListPermissionsByRoleInternalError() {
	roleID := uuid.New()

	req := &permissionV1.ListPermissionsByRoleRequest{
		Value: &commonV1.GetIdentifier{
			Identifier: &commonV1.GetIdentifier_Id{
				Id: roleID.String(),
			},
		},
	}

	s.permissionService.On("ListPermissionsByRole", mock.Anything, mock.Anything).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.ListPermissionsByRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.permissionService.AssertExpectations(s.T())
}
