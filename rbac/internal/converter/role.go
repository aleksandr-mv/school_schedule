package converter

import (
	"github.com/aleksandr-mv/school_schedule/rbac/internal/model"
	commonV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/common/v1"
	roleV1 "github.com/aleksandr-mv/school_schedule/shared/pkg/proto/role/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateRoleToDomain преобразует protobuf запрос в доменную модель создания роли
func CreateRoleToDomain(req *roleV1.CreateRoleRequest) *model.CreateRole {
	return &model.CreateRole{
		Name:        req.Name,
		Description: req.Description,
	}
}

// UpdateRoleToDomain преобразует protobuf запрос в доменную модель обновления роли
func UpdateRoleToDomain(req *roleV1.UpdateRoleRequest) (*model.UpdateRole, error) {
	updateRole := &model.UpdateRole{
		ID:          req.RoleId,
		Name:        req.Name,
		Description: req.Description,
	}

	return updateRole, nil
}

// ParseNameFilter извлекает фильтр имени из protobuf запроса
func ParseNameFilter(req *roleV1.ListRolesRequest) string {
	if req.NameFilter != nil {
		return *req.NameFilter
	}
	return ""
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
