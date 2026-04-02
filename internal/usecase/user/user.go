package user

import (
	"adoptme/internal/entity"
	"adoptme/internal/repo"
	"context"
	"fmt"
)

type UseCase struct {
	shelterRepo   repo.ShelterRepo
	volunteerRepo repo.VolunteerRepo
}

func New(shRepo repo.ShelterRepo, vlRepo repo.VolunteerRepo) *UseCase {
	return &UseCase{
		shelterRepo:   shRepo,
		volunteerRepo: vlRepo,
	}
}

func (uc *UseCase) RegisterShelter(ctx context.Context, sh entity.Shelter) error {
	err := uc.shelterRepo.Create(ctx, sh)
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterShelter - uc.repo.Create: %w", err)
	}

	return nil
}

func (uc *UseCase) RegisterVolunteer(ctx context.Context, v entity.Volunteer) error {
	err := uc.volunteerRepo.Create(ctx, v)
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterVolunteer - uc.repo.Create: %w", err)
	}

	return nil
}
