package entity

import "github.com/google/uuid"

type Volunteer struct {
	ID      uuid.UUID
	Email   string
	Name    string
	Surname string
}
