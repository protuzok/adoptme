package request

import "adoptme/internal/entity"

type RegisterAnimal struct {
	Name      string           `json:"name" validate:"required,min=2,max=255" example:"Rex"`
	OwnerID   string           `json:"owner_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	OwnerType entity.OwnerType `json:"owner_type" validate:"required" example:"shelter"`
} // @name v1.RegisterAnimal

type Transfer struct {
	NewOwnerID   string           `json:"new_owner_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	NewOwnerType entity.OwnerType `json:"new_owner_type" validate:"required" example:"shelter"`
} // @name v1.Transfer
