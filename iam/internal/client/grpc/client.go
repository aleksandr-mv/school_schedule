package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
)

type RBACClient interface {
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*model.RoleWithPermissions, error)
}
