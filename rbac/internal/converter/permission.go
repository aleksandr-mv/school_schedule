package converter

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
)

func PermissionsToProto(permissions []*model.Permission) []*commonV1.Permission {
	result := make([]*commonV1.Permission, len(permissions))
	for i, permission := range permissions {
		result[i] = &commonV1.Permission{
			Id:       permission.ID.String(),
			Resource: permission.Resource,
			Action:   permission.Action,
		}
	}
	return result
}
