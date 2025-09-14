package model

type EnrichedRoleRow struct {
	Role
	PermissionsJSON []byte `db:"permissions"`
}
