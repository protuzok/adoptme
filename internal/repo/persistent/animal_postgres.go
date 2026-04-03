package persistent

import (
	"adoptme/internal/entity"
	"adoptme/internal/repo"
	"adoptme/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AnimalRepo struct {
	*postgres.Postgres
}

func NewAnimal(pg *postgres.Postgres) *AnimalRepo {
	return &AnimalRepo{pg}
}

func (a AnimalRepo) Create(ctx context.Context, an entity.Animal) error {
	sql, args, err := a.Builder.
		Insert("animals").
		Columns("id, name, owner_id, owner_type").
		Values(an.ID, an.Name, an.OwnerID, an.OwnerType).
		ToSql()
	if err != nil {
		return fmt.Errorf("AnimalRepo - Create - a.Builder: %w", err)
	}

	_, err = a.Pool.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("AnimalRepo - Create - a.Pool.Exec: %w", err)
	}

	return nil
}

func (a AnimalRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Animal, error) {
	sql, args, err := a.Builder.
		Select("id, name, owner_id, owner_type").
		From("animals").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Animal{}, fmt.Errorf("AnimalRepo - GetByID - a.Builder: %w", err)
	}

	row := a.Pool.QueryRow(ctx, sql, args)

	var an entity.Animal
	err = row.Scan(&an.ID, &an.Name, &an.OwnerID, &an.OwnerType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Animal{}, fmt.Errorf("AnimalRepo - GetByID - row.Scan: %w", repo.ErrNotFound)
		}
		return entity.Animal{}, fmt.Errorf("AnimalRepo - GetByID - row.Scan: %w", err)
	}

	return an, nil
}

func (a AnimalRepo) UpdateOwner(ctx context.Context, animalID uuid.UUID, ownerID uuid.UUID, ownerType entity.OwnerType) error {
	sql, args, err := a.Builder.
		Update("animals").
		Set("owner_id", ownerID).
		Set("owner_type", ownerType).
		Where(squirrel.Eq{"id": animalID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("AnimalRepo - UpdateOwner - a.Builder: %w", err)
	}

	_, err = a.Pool.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("AnimalRepo - UpdateOwner - a.Pool.Exec: %w", err)
	}

	return nil
}
