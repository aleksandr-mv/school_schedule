package model

// CreateRole представляет данные для создания роли
type CreateRole struct {
	Name        string `validate:"required,min=2,max=50"`
	Description string `validate:"max=500"`
}

func (cr *CreateRole) Validate() error {
	return validate.Struct(cr)
}
