package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"adoptme/internal/entity"
	"adoptme/internal/usecase/adoption"
)

func adoptionUseCase(t *testing.T) (*adoption.UseCase, *MockAnimalRepo, *MockShelterRepo, *MockVolunteerRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)

	anRepo := NewMockAnimalRepo(mockCtl)
	shRepo := NewMockShelterRepo(mockCtl)
	vlRepo := NewMockVolunteerRepo(mockCtl)

	useCase := adoption.New(anRepo, shRepo, vlRepo)

	return useCase, anRepo, shRepo, vlRepo
}

func TestRegisterAnimal(t *testing.T) {
	t.Parallel()

	uc, anRepo, shRepo, vlRepo := adoptionUseCase(t)

	ownerID := "00000000-0000-0000-0000-000000000001"

	type args struct {
		name      string
		ownerID   string
		ownerType entity.OwnerType
	}
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "normal register (shelter owner)",
			args: args{name: "Borys", ownerID: ownerID, ownerType: entity.OwnerTypeShelter},
			mock: func() {
				shRepo.EXPECT().GetByID(gomock.Any(), ownerID).Return(entity.Shelter{ID: ownerID}, nil)
				anRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Animal{})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "normal register (volunteer owner)",
			args: args{name: "Murchyk", ownerID: ownerID, ownerType: entity.OwnerTypeVolunteer},
			mock: func() {
				vlRepo.EXPECT().GetByID(gomock.Any(), ownerID).Return(entity.Volunteer{ID: ownerID}, nil)
				anRepo.EXPECT().Store(gomock.Any(), gomock.AssignableToTypeOf(&entity.Animal{})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error: shelter not found",
			args: args{name: "Borys", ownerID: ownerID, ownerType: entity.OwnerTypeShelter},
			mock: func() {
				shRepo.EXPECT().GetByID(gomock.Any(), ownerID).Return(entity.Shelter{}, entity.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "error: invalid owner type",
			args: args{name: "Borys", ownerID: ownerID, ownerType: "alien"},
			mock: func() {
				// No mock calls expected
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			_, err := uc.RegisterAnimal(context.Background(), localTc.args.name, localTc.args.ownerID, localTc.args.ownerType)

			if localTc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTransferAnimal(t *testing.T) {
	t.Parallel()

	uc, anRepo, shRepo, vlRepo := adoptionUseCase(t)

	animalID := "00000000-0000-0000-0000-000000000002"
	newOwnerID := "00000000-0000-0000-0000-000000000003"

	type args struct {
		animalID     string
		newOwnerID   string
		newOwnerType entity.OwnerType
	}
	tests := []struct {
		name    string
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "normal transfer (to volunteer)",
			args: args{animalID: animalID, newOwnerID: newOwnerID, newOwnerType: entity.OwnerTypeVolunteer},
			mock: func() {
				// Check if the animal exists
				anRepo.EXPECT().GetByID(gomock.Any(), animalID).Return(entity.Animal{ID: animalID}, nil)
				// Check if the volunteer (new owner) exists
				vlRepo.EXPECT().GetByID(gomock.Any(), newOwnerID).Return(entity.Volunteer{ID: newOwnerID}, nil)
				// Update the owner
				anRepo.EXPECT().UpdateOwner(gomock.Any(), animalID, newOwnerID, entity.OwnerTypeVolunteer).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "normal transfer (to shelter)",
			args: args{animalID: animalID, newOwnerID: newOwnerID, newOwnerType: entity.OwnerTypeShelter},
			mock: func() {
				anRepo.EXPECT().GetByID(gomock.Any(), animalID).Return(entity.Animal{ID: animalID}, nil)
				shRepo.EXPECT().GetByID(gomock.Any(), newOwnerID).Return(entity.Shelter{ID: newOwnerID}, nil)
				anRepo.EXPECT().UpdateOwner(gomock.Any(), animalID, newOwnerID, entity.OwnerTypeShelter).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error: animal not found",
			args: args{animalID: animalID, newOwnerID: newOwnerID, newOwnerType: entity.OwnerTypeVolunteer},
			mock: func() {
				anRepo.EXPECT().GetByID(gomock.Any(), animalID).Return(entity.Animal{}, entity.ErrAnimalNotFound)
			},
			wantErr: true,
		},
		{
			name: "error: new owner (shelter) not found",
			args: args{animalID: animalID, newOwnerID: newOwnerID, newOwnerType: entity.OwnerTypeShelter},
			mock: func() {
				anRepo.EXPECT().GetByID(gomock.Any(), animalID).Return(entity.Animal{ID: animalID}, nil)
				shRepo.EXPECT().GetByID(gomock.Any(), newOwnerID).Return(entity.Shelter{}, entity.ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		localTc := tc
		t.Run(localTc.name, func(t *testing.T) {
			localTc.mock()

			err := uc.TransferAnimal(context.Background(), localTc.args.animalID, localTc.args.newOwnerID, localTc.args.newOwnerType)

			if localTc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
