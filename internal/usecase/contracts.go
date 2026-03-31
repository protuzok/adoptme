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
		ListShelters()
		ListVolunteer()
	}

	Adoption interface {
		RegisterAnimal()
		TransferAnimal()
	}
)
