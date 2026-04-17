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

const _defaultShelterCap = 64

type ShelterRepo struct {
	*postgres.Postgres
}

func NewShelterRepo(pg *postgres.Postgres) *ShelterRepo {
	return &ShelterRepo{pg}
}

func (r ShelterRepo) Store(ctx context.Context, sh *entity.Shelter) error {
	sql, args, err := r.Builder.
		Insert("shelters").
		Columns("id, email, name, password_hash, created_at, updated_at").
		Values(sh.ID, sh.Email, sh.Name, sh.PasswordHash, sh.CreatedAt, sh.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("ShelterRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}

		return fmt.Errorf("ShelterRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r ShelterRepo) GetByID(ctx context.Context, id string) (entity.Shelter, error) {
	return r.getUser(ctx, "id", id)
}

func (r ShelterRepo) GetByEmail(ctx context.Context, email string) (entity.Shelter, error) {
	return r.getUser(ctx, "email", email)
}

func (r ShelterRepo) getUser(ctx context.Context, column, value string) (entity.Shelter, error) {
	sql, args, err := r.Builder.
		Select("id, email, name, password_hash, created_at, updated_at").
		From("shelters").
		Where(sq.Eq{column: value}).
		ToSql()
	if err != nil {
		return entity.Shelter{}, fmt.Errorf("ShelterRepo - getUser - r.Builder: %w", err)
	}

	sh := entity.Shelter{}

	err = r.Pool.QueryRow(ctx, sql, args...).
		Scan(&sh.ID, &sh.Email, &sh.Name, &sh.PasswordHash, &sh.CreatedAt, &sh.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Shelter{}, entity.ErrUserNotFound
		}
		return entity.Shelter{}, fmt.Errorf("ShelterRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return sh, nil
}

func (r ShelterRepo) List(ctx context.Context) ([]entity.Shelter, error) {
	sql, args, err := r.Builder.
		Select("id, email, name, password_hash, created_at, updated_at").
		From("shelters").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ShelterRepo - List - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ShelterRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	shelters := make([]entity.Shelter, 0, _defaultShelterCap)

	for rows.Next() {
		var sh entity.Shelter

		err = rows.Scan(&sh.ID, &sh.Email, &sh.Name, &sh.PasswordHash, &sh.CreatedAt, &sh.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("ShelterRepo - List - rows.Scan: %w", err)
		}

		shelters = append(shelters, sh)
	}

	return shelters, nil
}
