package repository

import "github.com/CamiloLeonP/parking-radar/internal/app/domain"

//go:generate mockgen -source=./sensor_repository.go -destination=./../../test/shared/mocks/mock_sensor_repository.go -package=mockgen
type ISensorRepository interface {
	Create(sensor *domain.Sensor) error
	GetByID(id uint) (*domain.Sensor, error)
	ListByParkingLot(parkingLotID uint) ([]domain.Sensor, error)
	ListByEsp32DeviceID(esp32DeviceID uint64) ([]domain.Sensor, error)
	Update(sensor *domain.Sensor) error
	Delete(id uint) error
}
