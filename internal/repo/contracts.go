package repo

import (
	"adoptme/internal/entity"
	"context"
)

type (
	ShelterRepo interface {
		Store(ctx context.Context, sh *entity.Shelter) error
		GetByID(ctx context.Context, id string) (entity.Shelter, error)
		GetByEmail(ctx context.Context, email string) (entity.Shelter, error)
		List(ctx context.Context) ([]entity.Shelter, error)
	}

	VolunteerRepo interface {
		Store(ctx context.Context, vl *entity.Volunteer) error
		GetByID(ctx context.Context, id string) (entity.Volunteer, error)
		GetByEmail(ctx context.Context, email string) (entity.Volunteer, error)
		List(ctx context.Context) ([]entity.Volunteer, error)
	}

	AnimalRepo interface {
		Store(ctx context.Context, an *entity.Animal) error
		GetByID(ctx context.Context, id string) (entity.Animal, error)
		UpdateOwner(ctx context.Context, animalID, ownerID string, ownerType entity.OwnerType) error
		List(ctx context.Context) ([]entity.Animal, error)
	}
)
