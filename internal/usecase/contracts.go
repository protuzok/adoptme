package usecase

import (
	"adoptme/internal/entity"
	"context"
)

type (
	Shelter interface {
		Register(ctx context.Context, name, email, password string) (entity.Shelter, error)
		Login(ctx context.Context, email, password string) (string, error)
	}

	Volunteer interface {
		Register(ctx context.Context, name, email, password string) (entity.Volunteer, error)
		Login(ctx context.Context, email, password string) (string, error)
	}

	Catalog interface {
		ListShelters(ctx context.Context) ([]entity.Shelter, error)
		ListVolunteer(ctx context.Context) ([]entity.Volunteer, error)
		ListAnimal(ctx context.Context) ([]entity.Animal, error)
	}

	Adoption interface {
		RegisterAnimal(ctx context.Context, name, ownerId string, ownerType entity.OwnerType) (entity.Animal, error)
		TransferAnimal(ctx context.Context, animalID, newOwnerID string, newOwnerType entity.OwnerType) error
	}
)
