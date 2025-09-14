package grpc

import (
	"context"

	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	"github.com/google/uuid"
)

type RBACClient interface {
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*model.RoleWithPermissions, error)
}
