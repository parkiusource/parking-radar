package usecase

import (
	"errors"
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
)

type ISensorUseCase interface {
	CreateSensor(req CreateSensorRequest) error
	GetSensor(sensorID uint) (*SensorResponse, error)
	UpdateSensor(sensorID uint, req UpdateSensorRequest) error
	DeleteSensor(sensorID uint) error
	ListSensorsByParkingLot(parkingLotID uint) ([]SensorResponse, error)
	GetSensorByDeviceAndNumber(deviceIdentifier string, sensorNumber int) (*SensorResponse, error)
}

type SensorUseCase struct {
	SensorRepository      repository.ISensorRepository
	Esp32DeviceRepository repository.IEsp32DeviceRepository
}

type CreateSensorRequest struct {
	ParkingLotID     uint   `json:"parking_lot_id"`
	DeviceIdentifier string `json:"device_identifier"` // Dirección MAC
	SensorNumber     int    `json:"sensor_number"`
	Status           string `json:"status"`
}

type UpdateSensorRequest struct {
	Status           string `json:"status"`
	DeviceIdentifier string `json:"device_identifier"`
	SensorNumber     int    `json:"sensor_number"`
}

type SensorResponse struct {
	ID            uint   `json:"id"`
	ParkingLotID  uint   `json:"parking_lot_id"`
	Esp32DeviceID uint   `json:"esp32_device_id"`
	Status        string `json:"status"`
}

func NewSensorUseCase(sensorRepo repository.ISensorRepository, esp32DeviceRepo repository.IEsp32DeviceRepository) ISensorUseCase {
	return &SensorUseCase{
		SensorRepository:      sensorRepo,
		Esp32DeviceRepository: esp32DeviceRepo,
	}
}

func (uc *SensorUseCase) CreateSensor(req CreateSensorRequest) error {
	device, err := uc.Esp32DeviceRepository.GetByDeviceIdentifier(req.DeviceIdentifier)
	if err != nil {
		return errors.New("device not found")
	}

	sensor := domain.Sensor{
		ParkingLotID:     req.ParkingLotID,
		Esp32DeviceID:    uint(device.ID),
		Status:           req.Status,
		SensorNumber:     req.SensorNumber,
		DeviceIdentifier: req.DeviceIdentifier,
	}

	return uc.SensorRepository.Create(&sensor)
}

func (uc *SensorUseCase) GetSensor(sensorID uint) (*SensorResponse, error) {
	sensor, err := uc.SensorRepository.GetByID(sensorID)
	if err != nil {
		return nil, err
	}

	response := &SensorResponse{
		ID:            sensor.ID,
		ParkingLotID:  sensor.ParkingLotID,
		Esp32DeviceID: sensor.Esp32DeviceID,
		Status:        sensor.Status,
	}

	return response, nil
}

func (uc *SensorUseCase) UpdateSensor(sensorID uint, req UpdateSensorRequest) error {
	sensor, err := uc.SensorRepository.GetByID(sensorID)
	if err != nil {
		return err
	}

	sensor.Status = req.Status
	return uc.SensorRepository.Update(sensor)
}

func (uc *SensorUseCase) DeleteSensor(sensorID uint) error {
	return uc.SensorRepository.Delete(sensorID)
}

func (uc *SensorUseCase) ListSensorsByParkingLot(parkingLotID uint) ([]SensorResponse, error) {
	sensors, err := uc.SensorRepository.ListByParkingLot(parkingLotID)
	if err != nil {
		return nil, err
	}

	var response []SensorResponse
	for _, sensor := range sensors {
		response = append(response, SensorResponse{
			ID:            sensor.ID,
			ParkingLotID:  sensor.ParkingLotID,
			Esp32DeviceID: sensor.Esp32DeviceID,
			Status:        sensor.Status,
		})
	}

	return response, nil
}

func (uc *SensorUseCase) GetSensorByDeviceAndNumber(deviceIdentifier string, sensorNumber int) (*SensorResponse, error) {
	sensor, err := uc.SensorRepository.GetByDeviceAndNumber(deviceIdentifier, sensorNumber)
	if err != nil {
		return nil, err
	}

	response := &SensorResponse{
		ID:            sensor.ID,
		ParkingLotID:  sensor.ParkingLotID,
		Esp32DeviceID: sensor.Esp32DeviceID,
		Status:        sensor.Status,
	}

	return response, nil
}
