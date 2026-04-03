package usecase_test

import (
	"adoptme/internal/entity"
	"adoptme/internal/usecase/user"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errInternalServErr = errors.New("internal server error")

func userUseCase(t *testing.T) (*user.UseCase, *MockShelterRepo, *MockVolunteerRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	// We normally use defer mockCtl.Finish() in old gomock, but t.Cleanup is better or mockCtl is tied to t since go1.14
	defer mockCtl.Finish()

	shRepo := NewMockShelterRepo(mockCtl)
	vlRepo := NewMockVolunteerRepo(mockCtl)

	useCase := user.New(shRepo, vlRepo)

	return useCase, shRepo, vlRepo
}

func TestRegisterShelter(t *testing.T) {
	t.Parallel()

	userUseCase, shRepo, _ := userUseCase(t)

	type args struct{ sh entity.Shelter }
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr error
	}{
		{
			name: "normal register",
			args: args{sh: entity.Shelter{Email: "turykruslanz@gmail.com", Name: "Ruslan"}},
			mock: func() {
				shRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(entity.Shelter{})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "repo error: generic error",
			args: args{sh: entity.Shelter{Name: "Fail Shelter"}},
			mock: func() {
				shRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(entity.Shelter{})).Return(errInternalServErr)
			},
			wantErr: errInternalServErr,
		},
		{
			name: "negative: context timeout from repo",
			args: args{sh: entity.Shelter{Name: "Timeout Shelter"}},
			mock: func() {
				shRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(entity.Shelter{})).Return(context.DeadlineExceeded)
			},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range tests {
		// this row you can delete in go version start from 1.22
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			err := userUseCase.RegisterShelter(context.Background(), localTc.args.sh)

			require.ErrorIs(t, err, localTc.wantErr)
		})
	}
}

func TestRegisterVolunteer(t *testing.T) {
	t.Parallel()

	userCase, _, vlRepo := userUseCase(t)

	type args struct{ vl entity.Volunteer }
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr error
	}{
		{
			name: "normal register",
			args: args{vl: entity.Volunteer{Email: "volunteer@example.com", Name: "Ivan", Surname: "Ivanov"}},
			mock: func() {
				vlRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(entity.Volunteer{})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "repo error: generic error",
			args: args{vl: entity.Volunteer{Email: "fail@example.com"}},
			mock: func() {
				vlRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(entity.Volunteer{})).Return(errInternalServErr)
			},
			wantErr: errInternalServErr,
		},
		{
			name: "negative: context timeout from repo",
			args: args{vl: entity.Volunteer{Email: "timeout@example.com"}},
			mock: func() {
				vlRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(entity.Volunteer{})).Return(context.DeadlineExceeded)
			},
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range tests {
		localTc := tc

		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			err := userCase.RegisterVolunteer(context.Background(), localTc.args.vl)

			require.ErrorIs(t, err, localTc.wantErr)
		})
	}
}
