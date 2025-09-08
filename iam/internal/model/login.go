package model

type LoginCredentials struct {
	Login    string `validate:"required,min=3,max=50"`
	Password string `validate:"required,min=6"`
}

func (lc *LoginCredentials) Validate() error {
	return validate.Struct(lc)
}
