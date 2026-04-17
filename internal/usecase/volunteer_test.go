package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"adoptme/internal/entity"
	"adoptme/internal/usecase/volunteer"
	"adoptme/pkg/jwt"
)

func volunteerUseCase(t *testing.T) (*volunteer.UseCase, *MockVolunteerRepo, *jwt.Manager) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	vlRepo := NewMockVolunteerRepo(mockCtl)
	jwtManager := jwt.New("test-secret", time.Hour)

	useCase := volunteer.New(vlRepo, jwtManager)

	return useCase, vlRepo, jwtManager
}

func TestVolunteerRegister(t *testing.T) {
	t.Parallel()

	uc, vlRepo, _ := volunteerUseCase(t)

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
			args: args{name: "Ruslan", email: "volunteer@gmail.com", password: "password123"},
			mock: func() {
				vlRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Volunteer{})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error: generic error",
			args: args{name: "Fail Volunteer", email: "fail@example.com", password: "password"},
			mock: func() {
				vlRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Volunteer{})).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "negative: context timeout from repo",
			args: args{name: "Timeout Volunteer", email: "timeout@example.com", password: "password"},
			mock: func() {
				vlRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Volunteer{})).Return(context.DeadlineExceeded)
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

				require.Empty(t, user)
			} else {
				require.NoError(t, err)

				require.NotEmpty(t, user.ID)
				require.Equal(t, localTc.args.name, user.Name)
				require.Equal(t, localTc.args.email, user.Email)
				require.NotEmpty(t, user.PasswordHash)
				require.NotZero(t, user.CreatedAt)
				require.NotZero(t, user.UpdatedAt)
			}
		})
	}
}

func TestVolunteerLogin(t *testing.T) {
	t.Parallel()

	uc, vlRepo, _ := volunteerUseCase(t)

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
				vlRepo.EXPECT().GetByEmail(gomock.Any(), "success@example.com").Return(entity.Volunteer{ID: "123", PasswordHash: string(hash)}, nil)
			},
			wantErr: false,
		},
		{
			name: "error: email not found",
			args: args{email: "notfound@example.com", password: "password"},
			mock: func() {
				vlRepo.EXPECT().GetByEmail(gomock.Any(), "notfound@example.com").Return(entity.Volunteer{}, entity.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "error: wrong password",
			args: args{email: "wrongpass@example.com", password: "wrongpassword"},
			mock: func() {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				vlRepo.EXPECT().GetByEmail(gomock.Any(), "wrongpass@example.com").Return(entity.Volunteer{ID: "123", PasswordHash: string(hash)}, nil)
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
