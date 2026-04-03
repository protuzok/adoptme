package repo

import (
	"adoptme/internal/entity"
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type (
	ShelterRepo interface {
		Create(context.Context, entity.Shelter) error
		GetByID(context.Context, uuid.UUID) (entity.Shelter, error)
		GetArray(context.Context) ([]entity.Shelter, error)
	}

	VolunteerRepo interface {
		Create(context.Context, entity.Volunteer) error
		GetByID(context.Context, uuid.UUID) (entity.Volunteer, error)
		GetArray(context.Context) ([]entity.Volunteer, error)
	}

	AnimalRepo interface {
		Create(context.Context, entity.Animal) error
		GetByID(context.Context, uuid.UUID) (entity.Animal, error)
		UpdateOwner(ctx context.Context, animalID uuid.UUID, ownerID uuid.UUID, ownerType entity.OwnerType) error
	}
)
