package model

// RoleWithPermissions представляет роль с правами доступа
type RoleWithPermissions struct {
	Role        *Role
	Permissions []*Permission
}
