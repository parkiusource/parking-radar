package repository

import "github.com/CamiloLeonP/parking-radar/internal/app/domain"

type ISensorRepository interface {
	Create(sensor *domain.Sensor) error
	GetByID(id uint) (*domain.Sensor, error)
	ListByParkingLot(parkingLotID uint) ([]domain.Sensor, error)
	Update(sensor *domain.Sensor) error
	Delete(id uint) error
}
