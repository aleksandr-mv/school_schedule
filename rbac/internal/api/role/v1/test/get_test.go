package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	roleV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/role/v1"
)

func (s *APISuite) TestGetSuccess() {
	roleID := uuid.New().String()
	role := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}
	permissions := []*model.Permission{
		{ID: uuid.New(), Resource: "users", Action: "read"},
		{ID: uuid.New(), Resource: "users", Action: "write"},
	}

	enrichedRole := &model.EnrichedRole{
		Role:        *role,
		Permissions: permissions,
	}

	req := &roleV1.GetRequest{
		RoleId: roleID,
	}

	s.roleService.On("Get", mock.Anything, roleID).Return(enrichedRole, nil).Once()

	resp, err := s.api.Get(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.NotNil(s.T(), resp.Data)

	// Проверяем роль
	assert.Equal(s.T(), role.ID.String(), resp.Data.Role.Id)
	assert.Equal(s.T(), role.Name, resp.Data.Role.Name)
	assert.Equal(s.T(), role.Description, resp.Data.Role.Description)

	// Проверяем права
	assert.Len(s.T(), resp.Data.Permissions, 2)
	assert.Equal(s.T(), permissions[0].ID.String(), resp.Data.Permissions[0].Id)
	assert.Equal(s.T(), permissions[0].Resource, resp.Data.Permissions[0].Resource)
	assert.Equal(s.T(), permissions[0].Action, resp.Data.Permissions[0].Action)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetRoleNotFound() {
	roleID := uuid.New().String()

	req := &roleV1.GetRequest{
		RoleId: roleID,
	}

	s.roleService.On("Get", mock.Anything, roleID).Return(nil, model.ErrRoleNotFound).Once()

	resp, err := s.api.Get(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetInternalError() {
	roleID := uuid.New().String()

	req := &roleV1.GetRequest{
		RoleId: roleID,
	}

	s.roleService.On("Get", mock.Anything, roleID).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.Get(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestGetValidation_InvalidUUID() {
	invalidUUID := "invalid-uuid"

	req := &roleV1.GetRequest{
		RoleId: invalidUUID,
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value must be a valid UUID")
}

func (s *APISuite) TestGetValidation_ValidRequest() {
	validUUID := uuid.New().String()

	req := &roleV1.GetRequest{
		RoleId: validUUID,
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}
