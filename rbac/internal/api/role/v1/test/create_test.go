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

func (s *APISuite) TestCreateSuccess() {
	roleID := uuid.New()
	name := "admin"
	description := "Administrator role"

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: description,
	}

	s.roleService.On("Create", mock.Anything, name, description).Return(roleID, nil).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), roleID.String(), resp.RoleId)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateAlreadyExists() {
	name := "admin"
	description := "Administrator role"

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: description,
	}

	s.roleService.On("Create", mock.Anything, name, description).Return(uuid.Nil, model.ErrRoleAlreadyExists).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.AlreadyExists, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateInternalError() {
	name := "admin"
	description := "Administrator role"

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: description,
	}

	s.roleService.On("Create", mock.Anything, name, description).Return(uuid.Nil, model.ErrInternal).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)

	grpcErr, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.Internal, grpcErr.Code())

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateValidation_EmptyName() {
	req := &roleV1.CreateRequest{
		Name:        "",
		Description: "Test role",
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value length must be between 2 and 50 runes")
}

func (s *APISuite) TestCreateValidation_NameTooShort() {
	req := &roleV1.CreateRequest{
		Name:        "a",
		Description: "Test role",
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value length must be between 2 and 50 runes")
}

func (s *APISuite) TestCreateValidation_NameTooLong() {
	req := &roleV1.CreateRequest{
		Name:        strings.Repeat("a", 51),
		Description: "Test role",
	}

	err := req.Validate()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "value length must be between 2 and 50 runes")
}

func (s *APISuite) TestCreateValidation_NameMinLength() {
	req := &roleV1.CreateRequest{
		Name:        "ab",
		Description: "Test role",
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}

func (s *APISuite) TestCreateValidation_NameMaxLength() {
	req := &roleV1.CreateRequest{
		Name:        strings.Repeat("a", 50),
		Description: "Test role",
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}

func (s *APISuite) TestCreateValidation_EmptyDescription() {
	req := &roleV1.CreateRequest{
		Name:        "admin",
		Description: "",
	}

	err := req.Validate()
	assert.NoError(s.T(), err)
}

func (s *APISuite) TestCreateEdgeCase_EmptyDescription() {
	roleID := uuid.New()
	name := "admin"

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: "", // Пустое описание должно быть валидным
	}

	s.roleService.On("Create", mock.Anything, name, "").Return(roleID, nil).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), roleID.String(), resp.RoleId)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateEdgeCase_UnicodeName() {
	roleID := uuid.New()
	name := "админ" // Unicode символы

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: "Роль администратора",
	}

	s.roleService.On("Create", mock.Anything, name, "Роль администратора").Return(roleID, nil).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), roleID.String(), resp.RoleId)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateEdgeCase_SpecialCharactersInName() {
	roleID := uuid.New()
	name := "admin-role_1" // Специальные символы

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: "Role with special characters",
	}

	s.roleService.On("Create", mock.Anything, name, "Role with special characters").Return(roleID, nil).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), roleID.String(), resp.RoleId)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateEdgeCase_VeryLongDescription() {
	roleID := uuid.New()
	name := "admin"
	description := strings.Repeat("Very long description. ", 100) // Очень длинное описание

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: description,
	}

	s.roleService.On("Create", mock.Anything, name, description).Return(roleID, nil).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), roleID.String(), resp.RoleId)

	s.roleService.AssertExpectations(s.T())
}

func (s *APISuite) TestCreateEdgeCase_WhitespaceInName() {
	roleID := uuid.New()
	name := "  admin  " // Пробелы в начале и конце

	req := &roleV1.CreateRequest{
		Name:        name,
		Description: "Test role",
	}

	// Ожидаем, что пробелы будут сохранены как есть
	s.roleService.On("Create", mock.Anything, name, "Test role").Return(roleID, nil).Once()

	resp, err := s.api.Create(s.ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), roleID.String(), resp.RoleId)

	s.roleService.AssertExpectations(s.T())
}
