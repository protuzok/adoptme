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
				shRepo.EXPECT().Create(context.Background(), entity.Shelter{Email: "turykruslanz@gmail.com", Name: "Ruslan"}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "repo error",
			args: args{sh: entity.Shelter{}},
			mock: func() {
				shRepo.EXPECT().Create(context.Background(), entity.Shelter{}).Return(errInternalServErr)
			},
			wantErr: errInternalServErr,
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
	// todo написати TestRegisterVolunteer
}
