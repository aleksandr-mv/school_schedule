package interceptor

import (
	"encoding/json"
	"strings"

	commonV1 "github.com/Alexander-Mandzhiev/school_schedule/shared/pkg/proto/common/v1"
)

// RoleParser отвечает за парсинг ролей, разрешений и notification methods из заголовков
type RoleParser struct{}

// NewRoleParser создает новый парсер заголовков
func NewRoleParser() *RoleParser {
	return &RoleParser{}
}

// ParseRolesWithPermissions создает структуры RoleWithPermissions из заголовков
func (p *RoleParser) ParseRolesWithPermissions(roleNames, permissionsList []string) []*commonV1.RoleWithPermissions {
	if len(roleNames) == 0 {
		return []*commonV1.RoleWithPermissions{}
	}

	permissionsMap := p.parsePermissionsMap(roleNames, permissionsList)
	result := p.buildRoleStructures(roleNames, permissionsMap)

	return result
}

// parsePermissionsMap парсит разрешения в map для быстрого поиска
func (p *RoleParser) parsePermissionsMap(roleNames, permissionsList []string) map[string][]*commonV1.Permission {
	permissionsMap := make(map[string][]*commonV1.Permission)

	for _, permStr := range permissionsList {
		parts := strings.Split(permStr, ":")
		if len(parts) != 2 {
			continue
		}

		resource := parts[0]
		action := parts[1]

		for _, roleName := range roleNames {
			if permissionsMap[roleName] == nil {
				permissionsMap[roleName] = make([]*commonV1.Permission, 0)
			}
			permissionsMap[roleName] = append(permissionsMap[roleName], &commonV1.Permission{
				Id:       "",
				Resource: resource,
				Action:   action,
			})
		}
	}

	return permissionsMap
}

// buildRoleStructures создает финальные RoleWithPermissions структуры
func (p *RoleParser) buildRoleStructures(roleNames []string, permissionsMap map[string][]*commonV1.Permission) []*commonV1.RoleWithPermissions {
	result := make([]*commonV1.RoleWithPermissions, 0, len(roleNames))
	for _, roleName := range roleNames {
		result = append(result, &commonV1.RoleWithPermissions{
			Role: &commonV1.Role{
				Id:          "",
				Name:        roleName,
				Description: "",
			},
			Permissions: permissionsMap[roleName],
		})
	}

	return result
}

// ParseNotificationMethods создает структуры NotificationMethod из JSON заголовка
func (p *RoleParser) ParseNotificationMethods(notificationJSON string) []*commonV1.NotificationMethod {
	if notificationJSON == "" {
		return nil
	}

	var tempNotifications []map[string]interface{}
	if err := json.Unmarshal([]byte(notificationJSON), &tempNotifications); err != nil {
		return nil
	}

	notificationMethods := make([]*commonV1.NotificationMethod, len(tempNotifications))
	for i, notif := range tempNotifications {
		notificationMethods[i] = &commonV1.NotificationMethod{
			ProviderName: getString(notif, "provider_name"),
			Target:       getString(notif, "target"),
		}
	}

	return notificationMethods
}

// getString безопасно извлекает строковое значение из map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
