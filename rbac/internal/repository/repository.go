package repository

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

// RoleRepository интерфейс для работы с ролями
type RoleRepository interface {
	Create(ctx context.Context, role *model.CreateRole) (*model.Role, error)
	Get(ctx context.Context, value string) (*model.Role, error) // value может быть UUID или имя роли
	Update(ctx context.Context, updateRole *model.UpdateRole) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, nameFilter string) ([]*model.Role, error)
}

// PermissionRepository интерфейс для работы с правами доступа
type PermissionRepository interface {
	Get(ctx context.Context, value string) (*model.Permission, error)
	List(ctx context.Context, filter *model.ListPermissionsFilter) ([]*model.Permission, error)
	GetByResourceAndAction(ctx context.Context, resource string, action string) (*model.Permission, error)
}

// UserRoleRepository интерфейс для работы с назначениями ролей пользователям
type UserRoleRepository interface {
	AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error
	RevokeRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error)
	GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor *string) ([]string, int32, *string, error)
	HasRole(ctx context.Context, userID, roleID string) (bool, error)
}

// RolePermissionRepository интерфейс для работы с назначениями прав ролям
type RolePermissionRepository interface {
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	ListPermissionsByRole(ctx context.Context, roleValue string) ([]*model.Permission, error)
	ListRolesByPermission(ctx context.Context, permissionValue string) ([]*model.Role, error)
	HasPermission(ctx context.Context, roleID, permissionID string) (bool, error)
}
