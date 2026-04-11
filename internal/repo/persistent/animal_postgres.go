package persistent

import (
	"adoptme/internal/entity"
	"adoptme/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const _defaultAnimalCap = 64

type AnimalRepo struct {
	*postgres.Postgres
}

func NewAnimal(pg *postgres.Postgres) *AnimalRepo {
	return &AnimalRepo{pg}
}

func (r AnimalRepo) Store(ctx context.Context, an *entity.Animal) error {
	sql, args, err := r.Builder.
		Insert("animals").
		Columns("id, owner_id, owner_type, name, created_at, updated_at").
		Values(an.ID, an.OwnerID, an.OwnerType, an.Name, an.CreatedAt, an.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("AnimalRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("AnimalRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r AnimalRepo) GetByID(ctx context.Context, id string) (entity.Animal, error) {
	sql, args, err := r.Builder.
		Select("id, owner_id, owner_type, name, created_at, updated_at").
		From("animals").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Animal{}, fmt.Errorf("AnimalRepo - GetByID - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args)

	var an entity.Animal
	err = row.Scan(&an.ID, &an.OwnerID, &an.OwnerType, &an.Name, &an.CreatedAt, &an.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Animal{}, entity.ErrAnimalNotFound
		}
		return entity.Animal{}, fmt.Errorf("AnimalRepo - GetByID - row.Scan: %w", err)
	}

	return an, nil
}

func (r AnimalRepo) UpdateOwner(ctx context.Context, animalID string, ownerID string, ownerType entity.OwnerType) error {
	sql, args, err := r.Builder.
		Update("animals").
		Set("owner_id", ownerID).
		Set("owner_type", ownerType).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": animalID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("AnimalRepo - UpdateOwner - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("AnimalRepo - UpdateOwner - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r AnimalRepo) List(ctx context.Context) ([]entity.Animal, error) {
	sql, args, err := r.Builder.
		Select("id, owner_id, owner_type, name, created_at, updated_at").
		From("animals").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AnimalRepo - List - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("AnimalRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	animals := make([]entity.Animal, 0, _defaultAnimalCap)

	for rows.Next() {
		var an entity.Animal

		err = rows.Scan(&an.ID, &an.OwnerID, &an.OwnerType, &an.Name, &an.CreatedAt, &an.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("AnimalRepo - List - rows.Scan: %w", err)
		}

		animals = append(animals, an)
	}

	return animals, nil
}
