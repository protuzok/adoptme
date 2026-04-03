package catalog

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

func (u UseCase) ListShelters(ctx context.Context) ([]entity.Shelter, error) {
	shelters, err := u.shelterRepo.GetArray(ctx)
	if err != nil {
		return nil, fmt.Errorf("CatalogUseCase - ListShelters - u.shelterRepo.GetArray: %w", err)
	}

	return shelters, nil
}

func (u UseCase) ListVolunteer(ctx context.Context) ([]entity.Volunteer, error) {
	volunteers, err := u.volunteerRepo.GetArray(ctx)
	if err != nil {
		return nil, fmt.Errorf("CatalogUseCase - ListVolunteer - u.volunteerRepo.GetArray: %w", err)
	}

	return volunteers, nil
}
