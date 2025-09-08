package model

type CreateUser struct {
	Login    string `validate:"required,min=3,max=50"`
	Email    string `validate:"required,email,max=255"`
	Password string `validate:"required,min=6"`
}

func (cu *CreateUser) Validate() error {
	return validate.Struct(cu)
}
