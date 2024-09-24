package db

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"gorm.io/gorm"
)

type SensorRepositoryImpl struct {
	DB *gorm.DB
}

func (r *SensorRepositoryImpl) Create(sensor *domain.Sensor) error {
	return r.DB.Create(sensor).Error
}

func (r *SensorRepositoryImpl) GetByID(id uint) (*domain.Sensor, error) {
	var sensor domain.Sensor

	if err := r.DB.First(&sensor, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (r *SensorRepositoryImpl) ListByParkingLot(parkingLotID uint) ([]domain.Sensor, error) {
	var sensors []domain.Sensor
	if err := r.DB.Where("parking_lot_id = ?", parkingLotID).Find(&sensors).Error; err != nil {
		return nil, err
	}
	return sensors, nil
}

func (r *SensorRepositoryImpl) ListByEsp32DeviceID(esp32DeviceID uint64) ([]domain.Sensor, error) {
	var sensors []domain.Sensor
	if err := r.DB.Where("esp32_device_id = ?", esp32DeviceID).Find(&sensors).Error; err != nil {
		return nil, err
	}
	return sensors, nil
}

func (r *SensorRepositoryImpl) GetByDeviceAndNumber(deviceIdentifier string, sensorNumber int) (*domain.Sensor, error) {
	var sensor domain.Sensor

	if err := r.DB.First(&sensor, "device_identifier = ? AND sensor_number = ?", deviceIdentifier, sensorNumber).Error; err != nil {
		return nil, err
	}

	return &sensor, nil
}

func (r *SensorRepositoryImpl) Update(sensor *domain.Sensor) error {
	return r.DB.Save(sensor).Error
}

func (r *SensorRepositoryImpl) Delete(id uint) error {
	return r.DB.Delete(&domain.Sensor{}, "id = ?", id).Error
}
