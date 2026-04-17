package catalog

import (
	"adoptme/internal/entity"
	"adoptme/internal/repo"
	"context"
	"fmt"
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

func list[T any](ctx context.Context, listFn func(ctx context.Context) ([]T, error), callerName string) ([]T, error) {
	items, err := listFn(ctx)
	if err != nil {
		return nil, fmt.Errorf("CatalogUseCase - %s - listFn: %w", callerName, err)
	}

	return items, nil
}

func (u UseCase) ListShelters(ctx context.Context) ([]entity.Shelter, error) {
	return list(ctx, u.shelterRepo.List, "ListShelters")
}

func (u UseCase) ListVolunteer(ctx context.Context) ([]entity.Volunteer, error) {
	return list(ctx, u.volunteerRepo.List, "ListVolunteer")
}

func (u UseCase) ListAnimal(ctx context.Context) ([]entity.Animal, error) {
	return list(ctx, u.animalRepo.List, "ListAnimal")
}
