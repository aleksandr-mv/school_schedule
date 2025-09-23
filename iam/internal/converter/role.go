package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	commonV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/common/v1"
)

// RoleToProto конвертирует доменную модель роли в protobuf
func RoleToProto(r *model.Role) *commonV1.Role {
	if r == nil {
		return nil
	}

	role := &commonV1.Role{
		Id:          r.ID.String(),
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   timestamppb.New(r.CreatedAt),
	}

	if r.UpdatedAt != nil {
		role.UpdatedAt = timestamppb.New(*r.UpdatedAt)
	}

	return role
}

// PermissionToProto конвертирует доменную модель права в protobuf
func PermissionToProto(p *model.Permission) *commonV1.Permission {
	if p == nil {
		return nil
	}

	return &commonV1.Permission{
		Id:       p.ID.String(),
		Resource: p.Resource,
		Action:   p.Action,
	}
}

// RoleWithPermissionsToProto конвертирует доменную модель роли с правами в protobuf
func RoleWithPermissionsToProto(rwp *model.RoleWithPermissions) *commonV1.RoleWithPermissions {
	if rwp == nil {
		return nil
	}

	permissions := make([]*commonV1.Permission, 0, len(rwp.Permissions))
	for _, permission := range rwp.Permissions {
		permissions = append(permissions, PermissionToProto(permission))
	}

	return &commonV1.RoleWithPermissions{
		Role:        RoleToProto(rwp.Role),
		Permissions: permissions,
	}
}

// RoleWithPermissionsSliceToProto конвертирует слайс ролей с правами в protobuf
func RoleWithPermissionsSliceToProto(rwps []*model.RoleWithPermissions) []*commonV1.RoleWithPermissions {
	if rwps == nil {
		return nil
	}

	result := make([]*commonV1.RoleWithPermissions, 0, len(rwps))
	for _, rwp := range rwps {
		result = append(result, RoleWithPermissionsToProto(rwp))
	}

	return result
}
