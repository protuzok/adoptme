package request

type RegisterShelter struct {
	Name     string `json:"name" validate:"required,min=3,max=255" example:"AdoptMe"`
	Email    string `json:"email" validate:"required,email,max=254" example:"adoptme@example.com"`
	Password string `json:"password" validate:"required,min=6,max=72" example:"secret123"`
} // @name v1.RegisterShelter

type LoginShelter struct {
	Email    string `json:"email" validate:"required,email,max=254" example:"adoptme@example.com"`
	Password string `json:"password" validate:"required,min=6,max=72" example:"secret123"`
} // @name v1.LoginShelter
