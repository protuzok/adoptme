package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"adoptme/internal/entity"
	"adoptme/internal/repo"
	"adoptme/internal/usecase/catalog"
)

func catalogUseCase(t *testing.T) (*catalog.UseCase, *MockShelterRepo, *MockVolunteerRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	shRepo := NewMockShelterRepo(mockCtl)
	vlRepo := NewMockVolunteerRepo(mockCtl)

	useCase := catalog.New(shRepo, vlRepo)

	return useCase, shRepo, vlRepo
}

func TestListShelters(t *testing.T) {
	t.Parallel()

	uc, shRepo, _ := catalogUseCase(t)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Shelter
		wantErr bool
	}{
		{
			name: "success: list shelters",
			mock: func() {
				shRepo.EXPECT().GetArray(gomock.Any()).Return([]entity.Shelter{
					{Name: "Best Shelter"},
					{Name: "Happy Paws"},
				}, nil)
			},
			want: []entity.Shelter{
				{Name: "Best Shelter"},
				{Name: "Happy Paws"},
			},
			wantErr: false,
		},
		{
			name: "error: repo failed",
			mock: func() {
				shRepo.EXPECT().GetArray(gomock.Any()).Return(nil, repo.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			got, err := uc.ListShelters(context.Background())

			if localTc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, localTc.want, got) // test that the array actually matches
			}
		})
	}
}

func TestListVolunteer(t *testing.T) {
	t.Parallel()

	uc, _, vlRepo := catalogUseCase(t)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Volunteer
		wantErr bool
	}{
		{
			name: "success: list volunteers",
			mock: func() {
				vlRepo.EXPECT().GetArray(gomock.Any()).Return([]entity.Volunteer{
					{Name: "Ruslan", Surname: "Turyk"},
				}, nil)
			},
			want: []entity.Volunteer{
				{Name: "Ruslan", Surname: "Turyk"},
			},
			wantErr: false,
		},
		{
			name: "error: repo failed",
			mock: func() {
				vlRepo.EXPECT().GetArray(gomock.Any()).Return(nil, repo.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			got, err := uc.ListVolunteer(context.Background())

			if localTc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, localTc.want, got) // test that the array actually matches
			}
		})
	}
}
