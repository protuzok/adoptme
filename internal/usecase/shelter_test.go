package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"adoptme/internal/entity"
	"adoptme/internal/usecase/shelter"
	"adoptme/pkg/jwt"
)

func shelterUseCase(t *testing.T) (*shelter.UseCase, *MockShelterRepo, *jwt.Manager) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	shRepo := NewMockShelterRepo(mockCtl)
	jwtManager := jwt.New("test-secret", time.Hour)
	useCase := shelter.New(shRepo, jwtManager)

	return useCase, shRepo, jwtManager
}

func TestShelterRegister(t *testing.T) {
	t.Parallel()

	uc, shRepo, _ := shelterUseCase(t)

	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "normal register",
			args: args{name: "Ruslan", email: "turykruslanz@gmail.com", password: "password123"},
			mock: func() {
				shRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Shelter{})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "duplicate register",
			args: args{name: "Ruslan", email: "turykruslanz@gmail.com", password: "password123"},
			mock: func() {
				shRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Shelter{})).Return(entity.ErrUserAlreadyExists)
			},
			wantErr: true,
		},
		{
			name: "repo error: generic error",
			args: args{name: "Fail Shelter", email: "fail@example.com", password: "password"},
			mock: func() {
				shRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Shelter{})).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "negative: context timeout from repo",
			args: args{name: "Timeout Shelter", email: "timeout@example.com", password: "password"},
			mock: func() {
				shRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Shelter{})).Return(context.DeadlineExceeded)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			user, err := uc.Register(context.Background(), localTc.args.name, localTc.args.email, localTc.args.password)

			if localTc.wantErr {
				require.Error(t, err)
				assert.Empty(t, user)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, user.ID)
				assert.Equal(t, localTc.args.name, user.Name)
				assert.Equal(t, localTc.args.email, user.Email)
				assert.NotEmpty(t, user.PasswordHash)
				assert.NotZero(t, user.CreatedAt)
				assert.NotZero(t, user.UpdatedAt)
			}
		})
	}
}

func TestShelterLogin(t *testing.T) {
	t.Parallel()

	uc, shRepo, _ := shelterUseCase(t)

	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "success: correct credentials",
			args: args{email: "success@example.com", password: "password123"},
			mock: func() {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				shRepo.EXPECT().GetByEmail(gomock.Any(), "success@example.com").Return(entity.Shelter{ID: "123", PasswordHash: string(hash)}, nil)
			},
			wantErr: false,
		},
		{
			name: "error: email not found",
			args: args{email: "notfound@example.com", password: "password"},
			mock: func() {
				shRepo.EXPECT().GetByEmail(gomock.Any(), "notfound@example.com").Return(entity.Shelter{}, entity.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "error: wrong password",
			args: args{email: "wrongpass@example.com", password: "wrongpassword"},
			mock: func() {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				shRepo.EXPECT().GetByEmail(gomock.Any(), "wrongpass@example.com").Return(entity.Shelter{ID: "123", PasswordHash: string(hash)}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			token, err := uc.Login(context.Background(), localTc.args.email, localTc.args.password)

			if localTc.wantErr {
				require.Error(t, err)
				require.Empty(t, token)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)
			}
		})
	}
}
