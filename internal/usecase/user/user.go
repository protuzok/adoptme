package user

import (
	"adoptme/internal/entity"
	"adoptme/internal/repo"
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (u *UseCase) RegisterShelter(ctx context.Context, sh entity.Shelter) error {
	var err error

	sh.ID, err = uuid.NewV7()
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterShelter - uuid.NewV7: %w", err)
	}

	err = u.shelterRepo.Create(ctx, sh)
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterShelter - u.repo.Create: %w", err)
	}

	return nil
}

func (u *UseCase) RegisterVolunteer(ctx context.Context, vl entity.Volunteer) error {
	var err error

	vl.ID, err = uuid.NewV7()
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterVolunteer - uuid.NewV7: %w", err)
	}

	err = u.volunteerRepo.Create(ctx, vl)
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterVolunteer - u.repo.Create: %w", err)
	}

	return nil
}
