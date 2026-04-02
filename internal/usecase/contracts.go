package usecase

import (
	"adoptme/internal/entity"
	"context"
)

type (
	User interface {
		RegisterShelter(context.Context, entity.Shelter) error
		RegisterVolunteer(context.Context, entity.Volunteer) error
	}

	Catalog interface {
		ListShelters(context.Context) ([]entity.Shelter, error)
		ListVolunteer(context.Context) ([]entity.Shelter, error)
	}

	Adoption interface {
		RegisterAnimal(context.Context, entity.Animal) error
		TransferAnimal(context.Context) error // ?
	}
)
