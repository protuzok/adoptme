package repo

import (
	"adoptme/internal/entity"
	"context"
)

type (
	ShelterRepo interface {
		Create(context.Context, entity.Shelter) error
		GetArray(context.Context) ([]entity.Shelter, error)
	}

	VolunteerRepo interface {
		Create(context.Context, entity.Volunteer) error
		GetArray(context.Context) ([]entity.Volunteer, error)
	}

	AnimalRepo interface {
		Create(context.Context, entity.Animal) error
		UpdateOwner(context.Context) // ?
	}
)
