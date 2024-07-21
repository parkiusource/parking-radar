package db

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"gorm.io/gorm"
)

type ParkingLotRepository struct {
	DB *gorm.DB
}

func (r *ParkingLotRepository) Create(parkingLot *domain.ParkingLot) error {
	return r.DB.Create(parkingLot).Error
}

func (r *ParkingLotRepository) GetByID(id uint) (*domain.ParkingLot, error) {
	var parkingLot domain.ParkingLot
	if err := r.DB.First(&parkingLot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &parkingLot, nil
}

func (r *ParkingLotRepository) Update(parkingLot *domain.ParkingLot) error {
	return r.DB.Save(parkingLot).Error
}

func (r *ParkingLotRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.ParkingLot{}, "id = ?", id).Error
}

func (r *ParkingLotRepository) List() ([]domain.ParkingLot, error) {
	var parlingLots []domain.ParkingLot
	if err := r.DB.Find(&parlingLots).Error; err != nil {
		return nil, err
	}

	return parlingLots, nil
}
