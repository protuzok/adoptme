package persistent

import (
	"adoptme/internal/entity"
	"adoptme/pkg/postgres"
	"context"
	"fmt"
)

const _defaultShelterCap = 64

type ShelterRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *ShelterRepo {
	return &ShelterRepo{pg}
}

func (s ShelterRepo) Create(ctx context.Context, sh entity.Shelter) error {
	sql, args, err := s.Builder.
		Insert("shelter").
		Columns("email, name").
		Values(sh.Email, sh.Name).
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

func (s ShelterRepo) GetArray(ctx context.Context) ([]entity.Shelter, error) {
	sql, args, err := s.Builder.Select("email, name").From("shelter").ToSql()
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

		err = rows.Scan(&sh.Email, &sh.Name)
		if err != nil {
			return nil, fmt.Errorf("ShelterRepo - GetArray - rows.Scan: %w", err)
		}

		shelters = append(shelters, sh)
	}

	return shelters, nil
}
