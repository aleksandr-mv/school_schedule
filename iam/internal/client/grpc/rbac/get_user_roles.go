package rbac

import (
	"context"

	converter "github.com/aleksandr-mv/school_schedule/iam/internal/client/converter/rbac"
	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	rbacV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/user_role/v1"
	"github.com/google/uuid"
)

func (c *client) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*model.RoleWithPermissions, error) {
	req := &rbacV1.GetUserRolesRequest{
		UserId: userID.String(),
	}

	res, err := c.generatedClient.GetUserRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	return converter.GetUserRolesResponseToDomain(res), nil
}
