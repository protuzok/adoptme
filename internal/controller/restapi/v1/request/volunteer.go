package request

type RegisterVolunteer struct {
	Name     string `json:"name" validate:"required,min=3,max=255" example:"Anton"`
	Email    string `json:"email" validate:"required,email,max=254" example:"Anton@example.com"`
	Password string `json:"password" validate:"required,min=6,max=72" example:"secret123"`
} // @name v1.RegisterVolunteer

type LoginVolunteer struct {
	Email    string `json:"email" validate:"required,email,max=254" example:"Anton@example.com"`
	Password string `json:"password" validate:"required,min=6,max=72" example:"secret123"`
} // @name v1.LoginVolunteer
