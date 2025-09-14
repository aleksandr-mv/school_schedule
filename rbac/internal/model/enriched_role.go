package model

type EnrichedRole struct {
	Role        Role
	Permissions []*Permission
}
