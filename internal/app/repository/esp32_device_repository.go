package repository

import "github.com/CamiloLeonP/parking-radar/internal/app/domain"

type IEsp32DeviceRepository interface {
	Create(device *domain.Esp32Device) error
	GetByID(id uint64) (*domain.Esp32Device, error)
	GetByDeviceIdentifier(identifier string) (*domain.Esp32Device, error)
	ListByDeviceIdentifier(identifier string) ([]domain.Esp32Device, error)
	ListAll() ([]domain.Esp32Device, error)
	Update(device *domain.Esp32Device) error
	Delete(id uint64) error
}
