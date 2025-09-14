package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
)

// UpdateRoleToDomain преобразует protobuf запрос в доменную модель обновления роли
func UpdateRoleToDomain(req *roleV1.UpdateRequest) (*model.UpdateRole, error) {
	updateRole := &model.UpdateRole{
		ID:          req.RoleId,
		Name:        req.Name,
		Description: req.Description,
	}

	return updateRole, nil
}

// RoleToProto преобразует модель роли в protobuf
func RoleToProto(role *model.Role) *commonV1.Role {
	var updatedAt *timestamppb.Timestamp
	if role.UpdatedAt != nil {
		updatedAt = timestamppb.New(*role.UpdatedAt)
	}

	return &commonV1.Role{
		Id:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   timestamppb.New(role.CreatedAt),
		UpdatedAt:   updatedAt,
	}
}

// RolesToProto преобразует массив моделей ролей в protobuf
func RolesToProto(roles []*model.Role) []*commonV1.Role {
	result := make([]*commonV1.Role, len(roles))
	for i, role := range roles {
		result[i] = RoleToProto(role)
	}
	return result
}

// EnrichedRoleToProto преобразует модель обогащенной роли в protobuf
func EnrichedRoleToProto(enrichedRole *model.EnrichedRole) *commonV1.RoleWithPermissions {
	return &commonV1.RoleWithPermissions{
		Role:        RoleToProto(&enrichedRole.Role),
		Permissions: PermissionsToProto(enrichedRole.Permissions),
	}
}

// EnrichedRolesToProto преобразует массив моделей обогащенных ролей в protobuf
func EnrichedRolesToProto(enrichedRoles []*model.EnrichedRole) []*commonV1.RoleWithPermissions {
	result := make([]*commonV1.RoleWithPermissions, len(enrichedRoles))
	for i, enrichedRole := range enrichedRoles {
		result[i] = EnrichedRoleToProto(enrichedRole)
	}
	return result
}
