package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
)

type RoleServiceInterface interface {
	Create(ctx context.Context, name, description string) (uuid.UUID, error)
	Get(ctx context.Context, id string) (*model.EnrichedRole, error)
	Update(ctx context.Context, updateRole *model.UpdateRole) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Role, error)
}

type PermissionServiceInterface interface {
	List(ctx context.Context) ([]*model.Permission, error)
}

type UserRoleServiceInterface interface {
	Assign(ctx context.Context, userID, roleID, assignedBy string) error
	Revoke(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*model.EnrichedRole, error)
	GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor string) ([]string, *string, error)
}

type RolePermissionServiceInterface interface {
	Assign(ctx context.Context, roleID, permissionID string) error
	Revoke(ctx context.Context, roleID, permissionID string) error
}
