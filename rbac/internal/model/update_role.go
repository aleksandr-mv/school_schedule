package model

// UpdateRole представляет данные для обновления роли
type UpdateRole struct {
	ID          string  `validate:"required,uuid"`
	Name        *string `validate:"omitempty,min=2,max=50"`
	Description *string `validate:"omitempty,max=500"`
}

func (u *UpdateRole) Validate() error {
	return validate.Struct(u)
}
