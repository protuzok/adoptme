package adoption

import (
	"context"
	"errors"
	"fmt"

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

func (u *UseCase) RegisterAnimal(ctx context.Context, an entity.Animal) error {
	var err error
	switch an.OwnerType {
	case entity.OwnerTypeShelter:
		_, err = u.shelterRepo.GetByID(ctx, an.OwnerID)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return fmt.Errorf("AdoptionUseCase - RegisterAnimal - shelter not found: %w", err)
			}
			return fmt.Errorf("AdoptionUseCase - RegisterAnimal - shelter DB error: %w", err)
		}
	case entity.OwnerTypeVolunteer:
		_, err = u.volunteerRepo.GetByID(ctx, an.OwnerID)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return fmt.Errorf("AdoptionUseCase - RegisterAnimal - volunteer not found: %w", err)
			}
			return fmt.Errorf("AdoptionUseCase - RegisterAnimal - volunteer DB error: %w", err)
		}
	default:
		return fmt.Errorf("AdoptionUseCase - RegisterAnimal - invalid owner type: %s", an.OwnerType)
	}

	an.ID, err = uuid.NewV7()
	if err != nil {
		return fmt.Errorf("AdoptionUseCase - RegisterAnimal - uuid.NewV7: %w", err)
	}

	err = u.animalRepo.Create(ctx, an)
	if err != nil {
		return fmt.Errorf("AdoptionUseCase - RegisterAnimal - u.animalRepo.Create: %w", err)
	}

	return nil
}

func (u *UseCase) TransferAnimal(ctx context.Context, animalID uuid.UUID, newOwnerID uuid.UUID, newOwnerType entity.OwnerType) error {
	_, err := u.animalRepo.GetByID(ctx, animalID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return fmt.Errorf("AdoptionUseCase - TransferAnimal - animal not found: %w", err)
		}
		return fmt.Errorf("AdoptionUseCase - TransferAnimal - animal DB error: %w", err)
	}

	switch newOwnerType {
	case entity.OwnerTypeShelter:
		_, err = u.shelterRepo.GetByID(ctx, newOwnerID)
		if err != nil {
			return fmt.Errorf("AdoptionUseCase - TransferAnimal - shelter DB error: %w", err)
		}
	case entity.OwnerTypeVolunteer:
		_, err = u.volunteerRepo.GetByID(ctx, newOwnerID)
		if err != nil {
			return fmt.Errorf("AdoptionUseCase - TransferAnimal - volunteer DB error: %w", err)
		}
	default:
		return fmt.Errorf("AdoptionUseCase - TransferAnimal - invalid owner type: %s", newOwnerType)
	}

	err = u.animalRepo.UpdateOwner(ctx, animalID, newOwnerID, newOwnerType)
	if err != nil {
		return fmt.Errorf("AdoptionUseCase - TransferAnimal - u.animalRepo.UpdateOwner: %w", err)
	}

	return nil
}
