package db

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"gorm.io/gorm"
)

type Esp32DeviceRepositoryImpl struct {
	DB *gorm.DB
}

func (r *Esp32DeviceRepositoryImpl) Create(device *domain.Esp32Device) error {
	return r.DB.Create(device).Error
}

func (r *Esp32DeviceRepositoryImpl) GetByID(id uint64) (*domain.Esp32Device, error) {
	var device domain.Esp32Device
	if err := r.DB.First(&device, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *Esp32DeviceRepositoryImpl) ListByDeviceIdentifier(identifier string) ([]domain.Esp32Device, error) {
	var devices []domain.Esp32Device
	if err := r.DB.Where("device_identifier = ?", identifier).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Esp32DeviceRepositoryImpl) Update(device *domain.Esp32Device) error {
	return r.DB.Save(device).Error
}

func (r *Esp32DeviceRepositoryImpl) Delete(id uint64) error {
	return r.DB.Delete(&domain.Esp32Device{}, "id = ?", id).Error
}
