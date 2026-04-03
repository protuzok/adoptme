package entity

import "github.com/google/uuid"

type OwnerType string

const (
	OwnerTypeShelter   OwnerType = "shelter"
	OwnerTypeVolunteer OwnerType = "volunteer"
)

type Animal struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID
	OwnerType OwnerType
}
