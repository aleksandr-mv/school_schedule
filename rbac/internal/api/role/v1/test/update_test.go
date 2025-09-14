package role_test

import (
	"strings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

func (s *APISuite) TestUpdateSuccess() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	req := &roleV1.UpdateRequest{
		RoleId:      roleID.String(),
		Name:        &name,
		Description: &description,
	}

	// Создаем ожидаемый объект обновления
	expectedUpdate := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleService.On("Update", mock.Anything, expectedUpdate).Return(nil).Once()

	resp, err := s.api.Update(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdateNotFound() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	req := &roleV1.UpdateRequest{
		RoleId:      roleID.String(),
		Name:        &name,
		Description: &description,
	}

	expectedUpdate := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleService.On("Update", mock.Anything, expectedUpdate).Return(model.ErrRoleNotFound).Once()

	resp, err := s.api.Update(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.NotFound, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdateInternalError() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated administrator role"

	req := &roleV1.UpdateRequest{
		RoleId:      roleID.String(),
		Name:        &name,
		Description: &description,
	}

	expectedUpdate := &model.UpdateRole{
		ID:          roleID.String(),
		Name:        &name,
		Description: &description,
	}

	s.roleService.On("Update", mock.Anything, expectedUpdate).Return(model.ErrInternal).Once()

	resp, err := s.api.Update(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestUpdateValidation_InvalidUUID() {
	req := &roleV1.UpdateRequest{
		RoleId: "invalid-uuid",
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value must be a valid UUID")
}

func (s *APISuite) TestUpdateValidation_NameTooShort() {
	roleID := uuid.New()
	name := "a"

	req := &roleV1.UpdateRequest{
		RoleId: roleID.String(),
		Name:   &name,
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value length must be between 2 and 50 runes")
}

func (s *APISuite) TestUpdateValidation_NameTooLong() {
	roleID := uuid.New()
	name := strings.Repeat("a", 51)

	req := &roleV1.UpdateRequest{
		RoleId: roleID.String(),
		Name:   &name,
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value length must be between 2 and 50 runes")
}

func (s *APISuite) TestUpdateValidation_ValidRequest() {
	roleID := uuid.New()
	name := "admin"
	description := "Updated role"

	req := &roleV1.UpdateRequest{
		RoleId:      roleID.String(),
		Name:        &name,
		Description: &description,
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}

func (s *APISuite) TestUpdateValidation_ValidRequestWithoutOptionalFields() {
	roleID := uuid.New()

	req := &roleV1.UpdateRequest{
		RoleId:      roleID.String(),
		Name:        nil,
		Description: nil,
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}
