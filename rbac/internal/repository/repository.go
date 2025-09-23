package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
)

type RoleRepository interface {
	Create(ctx context.Context, name, description string) (uuid.UUID, error)
	Get(ctx context.Context, id string) (*model.Role, error)
	Update(ctx context.Context, updateRole *model.UpdateRole) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*model.Role, error)
}

type PermissionRepository interface {
	List(ctx context.Context) ([]*model.Permission, error)
}

type UserRoleRepository interface {
	Assign(ctx context.Context, userID, roleID string, assignedBy *string) error
	Revoke(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetRoleUsers(ctx context.Context, roleID string, limit int32, cursor string) ([]string, *string, error)
}

type RolePermissionRepository interface {
	Assign(ctx context.Context, roleID, permissionID string) error
	Revoke(ctx context.Context, roleID, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]*model.Permission, error)
}

type EnrichedRoleRepository interface {
	Set(ctx context.Context, role *model.EnrichedRole, expiresAt time.Time) error
	Get(ctx context.Context, id string) (*model.EnrichedRole, error)
	Delete(ctx context.Context, id string) error
}
