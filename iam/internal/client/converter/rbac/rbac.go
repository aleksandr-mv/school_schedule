package converter

import (
	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	rbacV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/user_role/v1"
)

// GetUserRolesResponseToDomain конвертирует protobuf ответ в доменные модели
func GetUserRolesResponseToDomain(resp *rbacV1.GetUserRolesResponse) []*model.RoleWithPermissions {
	if resp == nil || resp.Data == nil {
		return nil
	}

	rolesWithPermissions := make([]*model.RoleWithPermissions, 0, len(resp.Data))
	for _, rwp := range resp.Data {
		if rwp == nil {
			continue
		}

		permissions := make([]*model.Permission, 0, len(rwp.Permissions))
		for _, permission := range rwp.Permissions {
			permissions = append(permissions, PermissionToDomain(permission))
		}

		rolesWithPermissions = append(rolesWithPermissions, &model.RoleWithPermissions{
			Role:        RoleToDomain(rwp.Role),
			Permissions: permissions,
		})
	}

	return rolesWithPermissions
}
