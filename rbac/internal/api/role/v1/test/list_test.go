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

func (s *APISuite) TestListSuccess() {
	role1 := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}
	role2 := &model.Role{
		ID:          uuid.New(),
		Name:        "user",
		Description: "User role",
		CreatedAt:   time.Now(),
	}

	expectedRoles := []*model.Role{role1, role2}

	req := &roleV1.ListRequest{}

	s.roleService.On("List", mock.Anything).Return(expectedRoles, nil).Once()

	resp, err := s.api.List(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 2)

	// Проверяем первую роль
	assert.Equal(s.T(), role1.ID.String(), resp.Data[0].Id)
	assert.Equal(s.T(), role1.Name, resp.Data[0].Name)
	assert.Equal(s.T(), role1.Description, resp.Data[0].Description)

	// Проверяем вторую роль
	assert.Equal(s.T(), role2.ID.String(), resp.Data[1].Id)
	assert.Equal(s.T(), role2.Name, resp.Data[1].Name)
	assert.Equal(s.T(), role2.Description, resp.Data[1].Description)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestListEmptyResult() {
	expectedRoles := []*model.Role{}
	req := &roleV1.ListRequest{}

	s.roleService.On("List", mock.Anything).Return(expectedRoles, nil).Once()

	resp, err := s.api.List(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Len(s.T(), resp.Data, 0)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestListInternalError() {
	req := &roleV1.ListRequest{}

	s.roleService.On("List", mock.Anything).Return(nil, model.ErrInternal).Once()

	resp, err := s.api.List(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestListValidation_ValidRequest() {
	req := &roleV1.ListRequest{}

	err := req.Validate()
	assert.NoError(s.T(), err)
}
