package db

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"gorm.io/gorm"
)

type ParkingLotRepositoryImpl struct {
	DB *gorm.DB
}

func (r *ParkingLotRepositoryImpl) Create(parkingLot *domain.ParkingLot) error {
	return r.DB.Create(parkingLot).Error
}

func (r *ParkingLotRepositoryImpl) GetByID(id uint) (*domain.ParkingLot, error) {
	var parkingLot domain.ParkingLot
	if err := r.DB.First(&parkingLot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &parkingLot, nil
}

func (r *ParkingLotRepositoryImpl) Update(parkingLot *domain.ParkingLot) error {
	return r.DB.Save(parkingLot).Error
}

func (r *ParkingLotRepositoryImpl) Delete(id uint) error {
	return r.DB.Delete(&domain.ParkingLot{}, "id = ?", id).Error
}

func (r *ParkingLotRepositoryImpl) List() ([]domain.ParkingLot, error) {
	var parkingLots []domain.ParkingLot
	if err := r.DB.Find(&parkingLots).Error; err != nil {
		return nil, err
	}

	return parkingLots, nil
}
