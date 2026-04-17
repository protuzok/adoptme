package request

type RegisterShelter struct {
	Name     string `json:"name" validate:"required,min=3,max=255" example:"johndoe"`
	Email    string `json:"email"    validate:"required,email"         example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6"         example:"secret123"`
} // @name v1.RegisterShelter

type LoginShelter struct {
	Email    string `json:"email"    validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required"       example:"secret123"`
} // @name v1.LoginShelter
