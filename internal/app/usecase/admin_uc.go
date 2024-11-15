package usecase

import (
	"errors"
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
)

//go:generate mockgen -source=./admin_uc.go -destination=./../../test/parking/mocks/mock_admin_uc.go -package=mocks
type IAdminUseCase interface {
	RegisterAdmin(adminID string) error
	CompleteAdminProfile(adminID string, profileData domain.AdminProfileData) error
	GetAdminProfile(adminID string) (*domain.Admin, error)
	GetParkingLotsByAdmin(adminUUID string) ([]domain.ParkingLot, error)
}

type AdminUseCase struct {
	AdminRepository      repository.IAdminRepository
	ParkingLotRepository repository.IParkingLotRepository
}

func (uc *AdminUseCase) GetParkingLotsByAdmin(adminUUID string) ([]domain.ParkingLot, error) {

	admin, err := uc.AdminRepository.FindByAuth0UUID(adminUUID)
	if err != nil {
		return nil, err
	}
	return uc.ParkingLotRepository.FindByAdminID(admin.ID)
}

func NewAdminUseCase(adminRepo repository.IAdminRepository) IAdminUseCase {
	return &AdminUseCase{
		AdminRepository: adminRepo,
	}
}

func (uc *AdminUseCase) RegisterAdmin(adminID string) error {
	exists, err := uc.AdminRepository.ExistsByAuth0UUID(adminID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("admin already registered")
	}

	admin := domain.Admin{
		Auth0UUID: adminID,
	}
	return uc.AdminRepository.Create(&admin)
}

func (uc *AdminUseCase) CompleteAdminProfile(adminID string, profileData domain.AdminProfileData) error {
	admin, err := uc.AdminRepository.FindByAuth0UUID(adminID)
	if err != nil {
		return err
	}

	admin.NIT = profileData.NIT
	admin.PhotoURL = profileData.PhotoURL
	admin.ContactPhone = profileData.ContactPhone

	return uc.AdminRepository.Update(admin)
}

func (uc *AdminUseCase) GetAdminProfile(adminID string) (*domain.Admin, error) {
	return uc.AdminRepository.FindByAuth0UUID(adminID)
}
