package model

type ListPermissionsFilter struct {
	RoleID   *string `validate:"omitempty,uuid"`
	Resource *string `validate:"omitempty,min=1,max=100"`
	Action   *string `validate:"omitempty,min=1,max=50"`
}

func (l *ListPermissionsFilter) Validate() error {
	return validate.Struct(l)
}
