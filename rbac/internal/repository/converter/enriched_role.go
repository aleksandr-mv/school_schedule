package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexander-Mandzhiev/school_schedule/rbac/internal/model"
	commonv1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/common/v1"
)

type EnrichedRoleCacheConverter struct{}

func NewEnrichedRoleCacheConverter() *EnrichedRoleCacheConverter {
	return &EnrichedRoleCacheConverter{}
}

func (c *EnrichedRoleCacheConverter) ToCache(enrichedRole *model.EnrichedRole) ([]byte, error) {
	if enrichedRole == nil {
		return nil, fmt.Errorf("enriched role cannot be nil")
	}

	pbRole := &commonv1.RoleWithPermissions{
		Role: &commonv1.Role{
			Id:          enrichedRole.Role.ID.String(),
			Name:        enrichedRole.Role.Name,
			Description: enrichedRole.Role.Description,
			CreatedAt:   timestamppb.New(enrichedRole.Role.CreatedAt),
			UpdatedAt: func() *timestamppb.Timestamp {
				if enrichedRole.Role.UpdatedAt != nil {
					return timestamppb.New(*enrichedRole.Role.UpdatedAt)
				}
				return nil
			}(),
		},
		Permissions: PermissionsToCache(enrichedRole.Permissions),
	}

	data, err := proto.Marshal(pbRole)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	return data, nil
}

func (c *EnrichedRoleCacheConverter) FromCache(data []byte) (*model.EnrichedRole, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	pbRole := &commonv1.RoleWithPermissions{}
	if err := proto.Unmarshal(data, pbRole); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	if pbRole.Role == nil {
		return nil, fmt.Errorf("role cannot be nil")
	}

	roleID, err := uuid.Parse(pbRole.Role.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}

	var updatedAt *time.Time
	if pbRole.Role.UpdatedAt != nil {
		t := pbRole.Role.UpdatedAt.AsTime()
		updatedAt = &t
	}

	permissions := PermissionsFromCache(pbRole.Permissions)

	enrichedRole := &model.EnrichedRole{
		Role: model.Role{
			ID:          roleID,
			Name:        pbRole.Role.Name,
			Description: pbRole.Role.Description,
			CreatedAt:   pbRole.Role.CreatedAt.AsTime(),
			UpdatedAt:   updatedAt,
		},
		Permissions: permissions,
	}

	return enrichedRole, nil
}

func PermissionsToCache(permissions []*model.Permission) []*commonv1.Permission {
	if permissions == nil {
		return nil
	}

	pbPermissions := make([]*commonv1.Permission, len(permissions))
	for i, p := range permissions {
		pbPermissions[i] = &commonv1.Permission{
			Id:       p.ID.String(),
			Resource: p.Resource,
			Action:   p.Action,
		}
	}
	return pbPermissions
}

func PermissionsFromCache(pbPermissions []*commonv1.Permission) []*model.Permission {
	if pbPermissions == nil {
		return nil
	}

	permissions := make([]*model.Permission, len(pbPermissions))
	for i, pbp := range pbPermissions {
		permissionID, err := uuid.Parse(pbp.Id)
		if err != nil {
			continue
		}

		permissions[i] = &model.Permission{
			ID:       permissionID,
			Resource: pbp.Resource,
			Action:   pbp.Action,
		}
	}
	return permissions
}
