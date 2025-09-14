package converter

import (
	"github.com/aleksandr-mv/school_schedule/iam/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
	"github.com/google/uuid"
)

// PermissionToDomain конвертирует protobuf Permission в доменную модель
func PermissionToDomain(p *commonV1.Permission) *model.Permission {
	if p == nil {
		return nil
	}

	permissionID, _ := uuid.Parse(p.Id)

	return &model.Permission{
		ID:       permissionID,
		Resource: p.Resource,
		Action:   p.Action,
	}
}
