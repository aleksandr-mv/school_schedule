package role_test

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
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

	// Проверяем кэш сначала
	s.enrichedRoleRepository.On("Get", mock.Anything, roleID).Return(nil, assert.AnError).Once()

	// Если нет в кэше, получаем из репозитория
	s.roleRepository.On("Get", mock.Anything, roleID).Return(role, nil).Once()
	s.rolePermissionRepository.On("GetRolePermissions", mock.Anything, roleID).Return(permissions, nil).Once()

	// Сохраняем в кэш
	s.enrichedRoleRepository.On("Set", mock.Anything, mock.MatchedBy(func(er *model.EnrichedRole) bool {
		return er.Role.ID == role.ID && len(er.Permissions) == 2
	}), mock.AnythingOfType("time.Time")).Return(nil).Once()

	result, err := s.service.Get(s.ctx, roleID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), role.ID, result.Role.ID)
	assert.Equal(s.T(), role.Name, result.Role.Name)
	assert.Equal(s.T(), role.Description, result.Role.Description)
	assert.Len(s.T(), result.Permissions, 2)
	assert.Equal(s.T(), permissions[0].Resource, result.Permissions[0].Resource)
	assert.Equal(s.T(), permissions[0].Action, result.Permissions[0].Action)

	s.roleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertExpectations(s.T())
	s.enrichedRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetFromCache() {
	roleID := uuid.New().String()
	role := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}
	permissions := []*model.Permission{
		{ID: uuid.New(), Resource: "users", Action: "read"},
	}

	expectedEnrichedRole := &model.EnrichedRole{
		Role:        *role,
		Permissions: permissions,
	}

	// Роль найдена в кэше
	s.enrichedRoleRepository.On("Get", mock.Anything, roleID).Return(expectedEnrichedRole, nil).Once()

	result, err := s.service.Get(s.ctx, roleID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), role.ID, result.Role.ID)
	assert.Equal(s.T(), role.Name, result.Role.Name)
	assert.Len(s.T(), result.Permissions, 1)

	// Репозитории не должны вызываться, если данные в кэше
	s.roleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertExpectations(s.T())
	s.enrichedRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRoleNotFound() {
	roleID := uuid.New().String()

	// Нет в кэше
	s.enrichedRoleRepository.On("Get", mock.Anything, roleID).Return(nil, assert.AnError).Once()

	// Роль не найдена в репозитории
	s.roleRepository.On("Get", mock.Anything, roleID).Return(nil, model.ErrRoleNotFound).Once()

	result, err := s.service.Get(s.ctx, roleID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrRoleNotFound, err)

	s.roleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertExpectations(s.T())
	s.enrichedRoleRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetPermissionsError() {
	roleID := uuid.New().String()
	role := &model.Role{
		ID:          uuid.New(),
		Name:        "admin",
		Description: "Administrator role",
		CreatedAt:   time.Now(),
	}

	// Нет в кэше
	s.enrichedRoleRepository.On("Get", mock.Anything, roleID).Return(nil, assert.AnError).Once()

	// Роль найдена, но ошибка при получении прав
	s.roleRepository.On("Get", mock.Anything, roleID).Return(role, nil).Once()
	s.rolePermissionRepository.On("GetRolePermissions", mock.Anything, roleID).Return(nil, model.ErrInternal).Once()

	result, err := s.service.Get(s.ctx, roleID)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), model.ErrInternal, err)

	s.roleRepository.AssertExpectations(s.T())
	s.rolePermissionRepository.AssertExpectations(s.T())
	s.enrichedRoleRepository.AssertExpectations(s.T())
}
