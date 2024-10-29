package db

import (
	"errors"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"gorm.io/gorm"
)

type ParkingLotRepositoryImpl struct {
	DB *gorm.DB
}

// Create adds a new parking lot to the database.
func (r *ParkingLotRepositoryImpl) Create(parkingLot *domain.ParkingLot) error {
	return r.DB.Create(parkingLot).Error
}

// GetByID retrieves a parking lot by its ID.
func (r *ParkingLotRepositoryImpl) GetByID(id uint) (*domain.ParkingLot, error) {
	var parkingLot domain.ParkingLot
	if err := r.DB.First(&parkingLot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &parkingLot, nil
}

// GetByIDWithAdmin retrieves a parking lot by ID and verifies ownership.
func (r *ParkingLotRepositoryImpl) GetByIDWithAdmin(parkingLotID uint, adminID string) (*domain.ParkingLot, error) {
	var parkingLot domain.ParkingLot
	if err := r.DB.Where("id = ? AND admin_id = ?", parkingLotID, adminID).First(&parkingLot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("forbidden: you don't have access to this parking lot")
		}
		return nil, err
	}
	return &parkingLot, nil
}

// Update modifies an existing parking lot.
func (r *ParkingLotRepositoryImpl) Update(parkingLot *domain.ParkingLot) error {
	return r.DB.Save(parkingLot).Error
}

// Delete removes a parking lot by its ID.
func (r *ParkingLotRepositoryImpl) Delete(id uint) error {
	return r.DB.Delete(&domain.ParkingLot{}, "id = ?", id).Error
}

// DeleteWithAdmin removes a parking lot after verifying ownership.
func (r *ParkingLotRepositoryImpl) DeleteWithAdmin(parkingLotID uint, adminID string) error {
	if _, err := r.GetByIDWithAdmin(parkingLotID, adminID); err != nil {
		return err
	}
	return r.DB.Delete(&domain.ParkingLot{}, "id = ?", parkingLotID).Error
}

// List retrieves all parking lots from the database.
func (r *ParkingLotRepositoryImpl) List() ([]domain.ParkingLot, error) {
	var parkingLots []domain.ParkingLot
	if err := r.DB.Find(&parkingLots).Error; err != nil {
		return nil, err
	}
	return parkingLots, nil
}
