package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (s *APISuite) TestListRolesSuccess() {
	role1 := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
	role2 := &model.Role{
		ID:          uuid.New(),
		Name:        "user",
		Description: "User role",
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
	expectedRoles := []*model.Role{role1, role2}

	req := &roleV1.ListRolesRequest{
		NameFilter: nil,
	}

	s.roleService.On("ListRoles", mock.Anything, mock.Anything).Return(expectedRoles, nil).Once()

	resp, err := s.api.ListRoles(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 2)

	for i, role := range resp.Data {
		assert.Equal(s.T(), expectedRoles[i].ID.String(), role.Id)
		assert.Equal(s.T(), expectedRoles[i].Name, role.Name)
		assert.Equal(s.T(), expectedRoles[i].Description, role.Description)
	}

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestListRolesSuccessWithFilter() {
	role1 := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
	expectedRoles := []*model.Role{role1}
	adminFilter := "admin"

	req := &roleV1.ListRolesRequest{
		NameFilter: &adminFilter,
	}

	s.roleService.On("ListRoles", mock.Anything, mock.Anything).Return(expectedRoles, nil).Once()

	resp, err := s.api.ListRoles(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 1)

	assert.Equal(s.T(), expectedRoles[0].ID.String(), resp.Data[0].Id)
	assert.Equal(s.T(), expectedRoles[0].Name, resp.Data[0].Name)
	assert.Equal(s.T(), expectedRoles[0].Description, resp.Data[0].Description)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestListRolesEmptyResult() {
	expectedRoles := []*model.Role{}
	nonexistentFilter := "nonexistent"

	req := &roleV1.ListRolesRequest{
		NameFilter: &nonexistentFilter,
	}

	s.roleService.On("ListRoles", mock.Anything, mock.Anything).Return(expectedRoles, nil).Once()

	resp, err := s.api.ListRoles(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 0)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestListRolesInternalError() {
	req := &roleV1.ListRolesRequest{
		NameFilter: nil,
	}

	s.roleService.On("ListRoles", mock.Anything, mock.Anything).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.ListRoles(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}
