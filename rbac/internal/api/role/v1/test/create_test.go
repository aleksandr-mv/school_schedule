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

func (s *APISuite) TestCreateRoleSuccess() {
	roleID := uuid.New()

	req := &roleV1.CreateRoleRequest{
		Name:        "admin",
		Description: "Administrator role",
	}

	s.roleService.On("CreateRole", mock.Anything, mock.Anything).Return(roleID, nil).Once()

	resp, err := s.api.CreateRole(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), roleID.String(), resp.RoleId)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateRoleValidationError() {
	req := &roleV1.CreateRoleRequest{
		Name:        "", // пустое имя
		Description: "Administrator role",
	}

	s.roleService.On("CreateRole", mock.Anything, mock.Anything).Return(uuid.Nil, model.ErrInvalidCredentials).Once()

	resp, err := s.api.CreateRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.InvalidArgument, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateRoleAlreadyExists() {
	req := &roleV1.CreateRoleRequest{
		Name:        "admin",
		Description: "Administrator role",
	}

	s.roleService.On("CreateRole", mock.Anything, mock.Anything).Return(uuid.Nil, model.ErrRoleAlreadyExists).Once()

	resp, err := s.api.CreateRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.AlreadyExists, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateRoleInternalError() {
	req := &roleV1.CreateRoleRequest{
		Name:        "admin",
		Description: "Administrator role",
	}

	s.roleService.On("CreateRole", mock.Anything, mock.Anything).Return(uuid.Nil, model.ErrInternal).Once()

	resp, err := s.api.CreateRole(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}
