package entity

import "time"

type OwnerType string

const (
	OwnerTypeShelter   OwnerType = "shelter"
	OwnerTypeVolunteer OwnerType = "volunteer"
)

type Animal struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OwnerID   string    `json:"owner_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OwnerType OwnerType `json:"owner_type" example:"volunteer"`
	Name      string    `json:"name" example:"Bobik"`
	CreatedAt time.Time `json:"created_at" example:"2026-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-01-01T00:00:00Z"`
}
