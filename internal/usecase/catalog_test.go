package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"adoptme/internal/entity"
	"adoptme/internal/usecase/catalog"
)

func catalogUseCase(t *testing.T) (*catalog.UseCase, *MockShelterRepo, *MockVolunteerRepo, *MockAnimalRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	shRepo := NewMockShelterRepo(mockCtl)
	vlRepo := NewMockVolunteerRepo(mockCtl)
	anRepo := NewMockAnimalRepo(mockCtl)

	useCase := catalog.New(anRepo, shRepo, vlRepo)

	return useCase, shRepo, vlRepo, anRepo
}

func TestListShelters(t *testing.T) {
	t.Parallel()

	uc, shRepo, _, _ := catalogUseCase(t)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Shelter
		wantErr bool
	}{
		{
			name: "success: list shelters",
			mock: func() {
				shRepo.EXPECT().List(gomock.Any()).Return([]entity.Shelter{
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
				shRepo.EXPECT().List(gomock.Any()).Return(nil, entity.ErrNotFound)
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

	uc, _, vlRepo, _ := catalogUseCase(t)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Volunteer
		wantErr bool
	}{
		{
			name: "success: list volunteers",
			mock: func() {
				vlRepo.EXPECT().List(gomock.Any()).Return([]entity.Volunteer{
					{Name: "Ruslan"},
				}, nil)
			},
			want: []entity.Volunteer{
				{Name: "Ruslan"},
			},
			wantErr: false,
		},
		{
			name: "error: repo failed",
			mock: func() {
				vlRepo.EXPECT().List(gomock.Any()).Return(nil, entity.ErrNotFound)
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

func TestListAnimal(t *testing.T) {
	t.Parallel()

	uc, _, _, anRepo := catalogUseCase(t)

	tests := []struct {
		name    string
		mock    func()
		want    []entity.Animal
		wantErr bool
	}{
		{
			name: "success: list animals",
			mock: func() {
				anRepo.EXPECT().List(gomock.Any()).Return([]entity.Animal{
					{Name: "Borys"},
				}, nil)
			},
			want: []entity.Animal{
				{Name: "Borys"},
			},
			wantErr: false,
		},
		{
			name: "error: repo failed",
			mock: func() {
				anRepo.EXPECT().List(gomock.Any()).Return(nil, entity.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			got, err := uc.ListAnimal(context.Background())

			if localTc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, localTc.want, got) // test that the array actually matches
			}
		})
	}
}
