package usecase

import (
	"adoptme/internal/entity"
	"context"

	"github.com/google/uuid"
)

type (
	User interface {
		RegisterShelter(context.Context, entity.Shelter) error
		RegisterVolunteer(context.Context, entity.Volunteer) error
	}

	Catalog interface {
		ListShelters(context.Context) ([]entity.Shelter, error)
		ListVolunteer(context.Context) ([]entity.Volunteer, error)
	}

	Adoption interface {
		RegisterAnimal(context.Context, entity.Animal) error
		TransferAnimal(ctx context.Context, animalID uuid.UUID, newOwnerID uuid.UUID, newOwnerType entity.OwnerType) error
	}
)
