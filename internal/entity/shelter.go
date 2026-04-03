package entity

import "github.com/google/uuid"

type Shelter struct {
	ID    uuid.UUID
	Email string
	Name  string
}
