package service

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	"github.com/google/uuid"
)

type RoleServiceInterface interface {
	CreateRole(ctx context.Context, createRole *model.CreateRole) (uuid.UUID, error)
	GetRole(ctx context.Context, value string) (*model.Role, error)
	UpdateRole(ctx context.Context, updateRole *model.UpdateRole) error
	DeleteRole(ctx context.Context, id string) error
	ListRoles(ctx context.Context, nameFilter string) ([]*model.Role, error)
}

type PermissionServiceInterface interface {
	GetPermission(ctx context.Context, value string) (*model.Permission, error)
	ListPermissions(ctx context.Context, filter *model.ListPermissionsFilter) ([]*model.Permission, error)
	CheckPermission(ctx context.Context, userID, resource, action string) error

	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	ListPermissionsByRole(ctx context.Context, roleValue string) ([]*model.Permission, error)
	ListRolesByPermission(ctx context.Context, permissionValue string) ([]*model.Role, error)
}

type UserRoleServiceInterface interface {
	AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error
	RevokeRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error)
	GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor *string) ([]string, int32, *string, error)
}
