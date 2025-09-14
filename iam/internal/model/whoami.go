package model

// WhoAMI представляет полную информацию о текущей сессии пользователя
type WhoAMI struct {
	Session              Session
	User                 User
	RolesWithPermissions []*RoleWithPermissions
}
