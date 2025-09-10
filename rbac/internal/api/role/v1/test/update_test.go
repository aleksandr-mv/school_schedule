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

func (s *APISuite) TestUpdateRoleSuccess() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	req := &roleV1.UpdateRoleRequest{
		RoleId:      roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleService.On("UpdateRole", mock.Anything, mock.Anything).Return(nil).Once()

	resp, err := s.api.UpdateRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdateRoleValidationError() {
	roleID := uuid.New()
	emptyName := ""
	description := "Updated administrator role"

	req := &roleV1.UpdateRoleRequest{
		RoleId:      roleID.String(),
		Name:        &emptyName, // пустое имя
		Description: &description,
	}

	s.roleService.On("UpdateRole", mock.Anything, mock.Anything).Return(model.ErrInvalidCredentials).Once()

	resp, err := s.api.UpdateRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.InvalidArgument, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdateRoleNotFound() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	req := &roleV1.UpdateRoleRequest{
		RoleId:      roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleService.On("UpdateRole", mock.Anything, mock.Anything).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.UpdateRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdateRoleInternalError() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	req := &roleV1.UpdateRoleRequest{
		RoleId:      roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleService.On("UpdateRole", mock.Anything, mock.Anything).Return(model.ErrInternal).Once()

	resp, err := s.api.UpdateRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}
