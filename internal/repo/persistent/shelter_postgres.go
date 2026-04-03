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

const _defaultShelterCap = 64

type ShelterRepo struct {
	*postgres.Postgres
}

func NewShelter(pg *postgres.Postgres) *ShelterRepo {
	return &ShelterRepo{pg}
}

func (s ShelterRepo) Create(ctx context.Context, sh entity.Shelter) error {
	sql, args, err := s.Builder.
		Insert("shelters").
		Columns("id, email, name").
		Values(sh.ID, sh.Email, sh.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("ShelterRepo - Create - s.Builder: %w", err)
	}

	_, err = s.Pool.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("ShelterRepo - Create - s.Pool.Exec: %w", err)
	}

	return nil
}

func (s ShelterRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Shelter, error) {
	sql, args, err := s.Builder.
		Select("id, email, name").
		From("shelters").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Shelter{}, fmt.Errorf("ShelterRepo - GetByID - s.Builder: %w", err)
	}

	row := s.Pool.QueryRow(ctx, sql, args)

	sh := entity.Shelter{}
	err = row.Scan(&sh.ID, &sh.Email, &sh.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Shelter{}, fmt.Errorf("ShelterRepo - GetByID - row.Scan: %w", repo.ErrNotFound)
		}
		return entity.Shelter{}, fmt.Errorf("ShelterRepo - GetByID - row.Scan: %w", err)
	}

	return sh, nil
}

func (s ShelterRepo) GetArray(ctx context.Context) ([]entity.Shelter, error) {
	sql, args, err := s.Builder.
		Select("id, email, name").
		From("shelters").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ShelterRepo - GetArray - s.Builder: %w", err)
	}

	rows, err := s.Pool.Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("ShelterRepo - GetArray - s.Pool.Query: %w", err)
	}
	defer rows.Close()

	shelters := make([]entity.Shelter, 0, _defaultShelterCap)

	for rows.Next() {
		sh := entity.Shelter{}

		err = rows.Scan(&sh.ID, &sh.Email, &sh.Name)
		if err != nil {
			return nil, fmt.Errorf("ShelterRepo - GetArray - rows.Scan: %w", err)
		}

		shelters = append(shelters, sh)
	}

	return shelters, nil
}
