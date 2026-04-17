package persistent

import (
	"adoptme/internal/entity"
	"adoptme/pkg/postgres"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const _defaultVolunteerCap = 64

type VolunteerRepo struct {
	*postgres.Postgres
}

func NewVolunteerRepo(pg *postgres.Postgres) *VolunteerRepo {
	return &VolunteerRepo{pg}
}

func (r VolunteerRepo) Store(ctx context.Context, vl *entity.Volunteer) error {
	sql, args, err := r.Builder.
		Insert("volunteers").
		Columns("id, email, name, password_hash, created_at, updated_at").
		Values(vl.ID, vl.Email, vl.Name, vl.PasswordHash, vl.CreatedAt, vl.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("VolunteerRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}

		return fmt.Errorf("VolunteerRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r VolunteerRepo) GetByID(ctx context.Context, id string) (entity.Volunteer, error) {
	return r.getUser(ctx, "id", id)
}

func (r VolunteerRepo) GetByEmail(ctx context.Context, email string) (entity.Volunteer, error) {
	return r.getUser(ctx, "email", email)
}

func (r VolunteerRepo) getUser(ctx context.Context, column, value string) (entity.Volunteer, error) {
	sql, args, err := r.Builder.
		Select("id, email, name, password_hash, created_at, updated_at").
		From("volunteers").
		Where(sq.Eq{column: value}).
		ToSql()
	if err != nil {
		return entity.Volunteer{}, fmt.Errorf("VolunteerRepo - getUser - r.Builder: %w", err)
	}

	vl := entity.Volunteer{}

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&vl.ID, &vl.Email, &vl.Name, &vl.PasswordHash, &vl.CreatedAt, &vl.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Volunteer{}, entity.ErrUserNotFound
		}
		return entity.Volunteer{}, fmt.Errorf("VolunteerRepo - getUser - r.Pool.QueryRow: %w", err)
	}

	return vl, nil
}

func (r VolunteerRepo) List(ctx context.Context) ([]entity.Volunteer, error) {
	sql, args, err := r.Builder.
		Select("id, email, name, password_hash, created_at, updated_at").
		From("volunteers").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("VolunteerRepo - List - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("VolunteerRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	volunteers := make([]entity.Volunteer, 0, _defaultVolunteerCap)

	for rows.Next() {
		var vl entity.Volunteer

		err = rows.Scan(&vl.ID, &vl.Email, &vl.Name, &vl.PasswordHash, &vl.CreatedAt, &vl.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("VolunteerRepo - List - rows.Scan: %w", err)
		}

		volunteers = append(volunteers, vl)
	}

	return volunteers, nil
}
