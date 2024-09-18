package repository

import "github.com/CamiloLeonP/parking-radar/internal/app/domain"

type IParkingLotRepository interface {
	Create(parkingLot *domain.ParkingLot) error
	GetByID(id uint) (*domain.ParkingLot, error)
	Update(parkingLot *domain.ParkingLot) error
	Delete(id uint) error
	List() ([]domain.ParkingLot, error)
}
