package converter

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
	permissionV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/permission/v1"
)

// ListPermissionsToDomain преобразует protobuf запрос в доменную модель фильтра разрешений
func ListPermissionsToDomain(req *permissionV1.ListPermissionsRequest) *model.ListPermissionsFilter {
	return &model.ListPermissionsFilter{
		RoleID:   req.RoleId,
		Resource: req.ResourceFilter,
		Action:   req.ActionFilter,
	}
}

func PermissionToProto(permission *model.Permission) *commonV1.Permission {
	return &commonV1.Permission{
		Id:       permission.ID.String(),
		Resource: permission.Resource,
		Action:   permission.Action,
	}
}

func PermissionsToProto(permissions []*model.Permission) []*commonV1.Permission {
	result := make([]*commonV1.Permission, 0, len(permissions))
	for _, permission := range permissions {
		result = append(result, PermissionToProto(permission))
	}
	return result
}

// PermissionsToListResponse преобразует массив прав в готовый ответ ListPermissionsResponse
func PermissionsToListResponse(permissions []*model.Permission) *permissionV1.ListPermissionsResponse {
	return &permissionV1.ListPermissionsResponse{
		Data: PermissionsToProto(permissions),
	}
}

// PermissionsToListByRoleResponse преобразует массив прав в готовый ответ ListPermissionsByRoleResponse
func PermissionsToListByRoleResponse(permissions []*model.Permission) *permissionV1.ListPermissionsByRoleResponse {
	return &permissionV1.ListPermissionsByRoleResponse{
		Data: PermissionsToProto(permissions),
	}
}
