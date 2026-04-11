package shelter

import (
	"adoptme/internal/entity"
	"adoptme/internal/repo"
	"adoptme/pkg/jwt"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	shRepo repo.ShelterRepo
	jwt    *jwt.Manager
}

func New(shRepo repo.ShelterRepo, jwt *jwt.Manager) *UseCase {
	return &UseCase{
		shRepo: shRepo,
		jwt:    jwt,
	}
}

func (u UseCase) Register(ctx context.Context, name, email, password string) (entity.Shelter, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return entity.Shelter{}, fmt.Errorf("ShelterUseCase - Register - uuid.NewV7: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.Shelter{}, fmt.Errorf("ShelterUseCase - Register - bcrypt.GenerateFromPassword: %w", err)
	}

	now := time.Now().UTC()

	user := entity.Shelter{
		ID:           id.String(),
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = u.shRepo.Store(ctx, &user)
	if err != nil {
		return entity.Shelter{}, fmt.Errorf("ShelterUseCase - Register - u.shRepo.Store: %w", err)
	}

	return user, nil
}

func (u UseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.shRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", entity.ErrInvalidCredentials
	}

	token, err := u.jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("ShelterUseCase - Login - u.jwt.GenerateToken: %w", err)
	}

	return token, nil
}
