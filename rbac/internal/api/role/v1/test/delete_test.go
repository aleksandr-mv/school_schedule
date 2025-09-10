package role_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (s *APISuite) TestDeleteRoleSuccess() {
	roleID := uuid.New()

	req := &roleV1.DeleteRoleRequest{
		RoleId: roleID.String(),
	}

	s.roleService.On("DeleteRole", mock.Anything, mock.Anything).Return(nil).Once()

	resp, err := s.api.DeleteRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestDeleteRoleNotFound() {
	roleID := uuid.New()

	req := &roleV1.DeleteRoleRequest{
		RoleId: roleID.String(),
	}

	s.roleService.On("DeleteRole", mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.DeleteRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	
	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestDeleteRoleInternalError() {
	roleID := uuid.New()

	req := &roleV1.DeleteRoleRequest{
		RoleId: roleID.String(),
	}

	s.roleService.On("DeleteRole", mock.Anything, mock.Anything).Return(model.ErrInternal).Once()

	resp, err := s.api.DeleteRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	
	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}