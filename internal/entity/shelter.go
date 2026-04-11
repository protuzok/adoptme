package entity

import "time"

type Shelter struct {
	ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email        string    `json:"email" example:"adoptme@example.com"`
	Name         string    `json:"name" example:"AdoptMe"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at" example:"2026-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2026-01-01T00:00:00Z"`
}
