package converter

import (
	"github.com/google/uuid"

	"github.com/Alexander-Mandzhiev/school_schedule/iam/internal/model"
	commonV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/common/v1"
)

// RoleToDomain конвертирует protobuf Role в доменную модель
func RoleToDomain(r *commonV1.Role) *model.Role {
	if r == nil {
		return nil
	}

	roleID, err := uuid.Parse(r.Id)
	if err != nil {
		return nil
	}

	role := &model.Role{
		ID:          roleID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt.AsTime(),
	}

	if r.UpdatedAt != nil {
		updatedAt := r.UpdatedAt.AsTime()
		role.UpdatedAt = &updatedAt
	}

	return role
}
