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

const _defaultVolunteerCap = 64

type VolunteerRepo struct {
	*postgres.Postgres
}

func NewVolunteer(pg *postgres.Postgres) *VolunteerRepo {
	return &VolunteerRepo{pg}
}

func (v VolunteerRepo) Create(ctx context.Context, vl entity.Volunteer) error {
	sql, args, err := v.Builder.
		Insert("volunteers").
		Columns("id, email, name, surname").
		Values(vl.ID, vl.Email, vl.Name, vl.Surname).
		ToSql()
	if err != nil {
		return fmt.Errorf("VolunteerRepo - Create - v.Builder: %w", err)
	}

	_, err = v.Pool.Exec(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("VolunteerRepo - Create - v.Pool.Exec: %w", err)
	}

	return nil
}

func (v VolunteerRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Volunteer, error) {
	sql, args, err := v.Builder.
		Select("id, email, name, surname").
		From("volunteers").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Volunteer{}, fmt.Errorf("VolunteerRepo - GetByID - v.Builder: %w", err)
	}

	row := v.Pool.QueryRow(ctx, sql, args)

	vl := entity.Volunteer{}
	err = row.Scan(&vl.ID, &vl.Email, &vl.Name, &vl.Surname)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Volunteer{}, fmt.Errorf("VolunteerRepo - GetByID - row.Scan: %w", repo.ErrNotFound)
		}
		return entity.Volunteer{}, fmt.Errorf("VolunteerRepo - GetByID - row.Scan: %w", err)
	}

	return vl, nil
}

func (v VolunteerRepo) GetArray(ctx context.Context) ([]entity.Volunteer, error) {
	sql, args, err := v.Builder.
		Select("id, email, name, surname").
		From("volunteers").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("VolunteerRepo - GetArray - v.Builder: %w", err)
	}

	rows, err := v.Pool.Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("VolunteerRepo - GetArray - v.Pool.Query: %w", err)
	}
	defer rows.Close()

	volunteers := make([]entity.Volunteer, 0, _defaultVolunteerCap)

	for rows.Next() {
		vl := entity.Volunteer{}

		err = rows.Scan(&vl.ID, &vl.Email, &vl.Name, &vl.Surname)
		if err != nil {
			return nil, fmt.Errorf("VolunteerRepo - GetArray - rows.Scan: %w", err)
		}

		volunteers = append(volunteers, vl)
	}

	return volunteers, nil
}
