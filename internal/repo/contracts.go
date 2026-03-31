package repo

import (
	"adoptme/internal/entity"
	"context"
)

type (
	UserRepo interface {
		CreateShelter(context.Context, entity.Shelter) error
		CreateVolunteer(context.Context, entity.Volunteer) error
	}

	CatalogRepo interface {
		GetShelter(context.Context) ([]entity.Shelter, error)
		GetVolunteer(context.Context) ([]entity.Volunteer, error)
	}

	AdoptionRepo interface {
		Create(context.Context, entity.Animal) error
		UpdateOwner(context.Context) // ?
	}
)
