package user

import (
	"adoptme/internal/entity"
	"adoptme/internal/repo"
	"context"
	"fmt"
)

type UseCase struct {
	repo repo.UserRepo
}

func New(r repo.UserRepo) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) RegisterShelter(ctx context.Context, sh entity.Shelter) error {
	err := uc.repo.CreateShelter(ctx, sh)
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterShelter - uc.repo.CreateShelter: %w", err)
	}

	return nil
}

func (uc *UseCase) RegisterVolunteer(ctx context.Context, v entity.Volunteer) error {
	err := uc.repo.CreateVolunteer(ctx, v)
	if err != nil {
		return fmt.Errorf("UserUseCase - RegisterVolunteer - uc.repo.CreateVolunteer: %w", err)
	}

	return nil
}
