package adoption

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"adoptme/internal/entity"
	"adoptme/internal/repo"
)

type UseCase struct {
	animalRepo    repo.AnimalRepo
	shelterRepo   repo.ShelterRepo
	volunteerRepo repo.VolunteerRepo
}

func New(anRepo repo.AnimalRepo, shRepo repo.ShelterRepo, vlRepo repo.VolunteerRepo) *UseCase {
	return &UseCase{
		animalRepo:    anRepo,
		shelterRepo:   shRepo,
		volunteerRepo: vlRepo,
	}
}

func (u *UseCase) RegisterAnimal(ctx context.Context, name, ownerId string, ownerType entity.OwnerType) (entity.Animal, error) {
	// Check if the owner exists
	var err error
	switch ownerType {
	case entity.OwnerTypeShelter:
		_, err = u.shelterRepo.GetByID(ctx, ownerId)
		if err != nil {
			if errors.Is(err, entity.ErrUserNotFound) {
				return entity.Animal{}, fmt.Errorf("AdoptionUseCase - RegisterAnimal: %w", err)
			}
			return entity.Animal{}, fmt.Errorf("AdoptionUseCase - RegisterAnimal - shelter DB error: %w", err)
		}
	case entity.OwnerTypeVolunteer:
		_, err = u.volunteerRepo.GetByID(ctx, ownerId)
		if err != nil {
			if errors.Is(err, entity.ErrUserNotFound) {
				return entity.Animal{}, fmt.Errorf("AdoptionUseCase - RegisterAnimal: %w", err)
			}
			return entity.Animal{}, fmt.Errorf("AdoptionUseCase - RegisterAnimal - volunteer DB error: %w", err)
		}
	default:
		return entity.Animal{}, fmt.Errorf("AdoptionUseCase - RegisterAnimal - invalid owner type: %s", ownerType)
	}

	// Fill fields for animal
	id, err := uuid.NewV7()
	if err != nil {
		return entity.Animal{}, fmt.Errorf("AdoptionUseCase - RegisterAnimal - uuid.NewV7: %w", err)
	}

	now := time.Now().UTC()

	animal := entity.Animal{
		ID:        id.String(),
		OwnerID:   ownerId,
		OwnerType: ownerType,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Store animal
	err = u.animalRepo.Store(ctx, &animal)
	if err != nil {
		return entity.Animal{}, fmt.Errorf("AdoptionUseCase - RegisterAnimal - u.animalRepo.Store: %w", err)
	}

	return animal, nil
}

func (u *UseCase) TransferAnimal(ctx context.Context, animalID, newOwnerID string, newOwnerType entity.OwnerType) error {
	// Check if the animal exists
	_, err := u.animalRepo.GetByID(ctx, animalID)
	if err != nil {
		if errors.Is(err, entity.ErrAnimalNotFound) {
			return fmt.Errorf("AdoptionUseCase - TransferAnimal: %w", err)
		}
		return fmt.Errorf("AdoptionUseCase - TransferAnimal - animal DB error: %w", err)
	}

	// Check if the owner exists
	switch newOwnerType {
	case entity.OwnerTypeShelter:
		_, err = u.shelterRepo.GetByID(ctx, newOwnerID)
		if err != nil {
			if errors.Is(err, entity.ErrUserNotFound) {
				return fmt.Errorf("AdoptionUseCase - TransferAnimal: %w", err)
			}
			return fmt.Errorf("AdoptionUseCase - TransferAnimal - shelter DB error: %w", err)
		}
	case entity.OwnerTypeVolunteer:
		_, err = u.volunteerRepo.GetByID(ctx, newOwnerID)
		if err != nil {
			if errors.Is(err, entity.ErrUserNotFound) {
				return fmt.Errorf("AdoptionUseCase - TransferAnimal: %w", err)
			}
			return fmt.Errorf("AdoptionUseCase - TransferAnimal - volunteer DB error: %w", err)
		}
	default:
		return fmt.Errorf("AdoptionUseCase - TransferAnimal - invalid owner type: %s", newOwnerType)
	}

	// Update owner
	err = u.animalRepo.UpdateOwner(ctx, animalID, newOwnerID, newOwnerType)
	if err != nil {
		return fmt.Errorf("AdoptionUseCase - TransferAnimal - u.animalRepo.UpdateOwner: %w", err)
	}

	return nil
}
