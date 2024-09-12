package usecase

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
)

type ISensorUseCase interface {
	CreateSensor(req CreateSensorRequest) error
	GetSensor(sensorID uint) (*SensorResponse, error) // Cambia el retorno al response struct
	UpdateSensor(sensorID uint, req UpdateSensorRequest) error
	DeleteSensor(sensorID uint) error
	ListSensorsByParkingLot(parkingLotID uint) ([]SensorResponse, error) // Cambia el retorno al response struct
}

type SensorUseCase struct {
	SensorRepository      repository.ISensorRepository
	Esp32DeviceRepository repository.IEsp32DeviceRepository
}

type CreateSensorRequest struct {
	ParkingLotID uint   `json:"parking_lot_id"`
	Status       string `json:"status"`
}

type UpdateSensorRequest struct {
	Status string `json:"status"`
}

type SensorResponse struct {
	ID           uint   `json:"id"`
	ParkingLotID uint   `json:"parking_lot_id"`
	Status       string `json:"status"`
}

func NewSensorUseCase(sensorRepo repository.ISensorRepository, esp32DeviceRepo repository.IEsp32DeviceRepository) ISensorUseCase {
	return &SensorUseCase{
		SensorRepository:      sensorRepo,
		Esp32DeviceRepository: esp32DeviceRepo,
	}
}

func (uc *SensorUseCase) CreateSensor(req CreateSensorRequest) error {
	sensor := domain.Sensor{
		ParkingLotID: req.ParkingLotID,
		Status:       req.Status,
	}

	return uc.SensorRepository.Create(&sensor)
}

func (uc *SensorUseCase) GetSensor(sensorID uint) (*SensorResponse, error) {
	sensor, err := uc.SensorRepository.GetByID(sensorID)
	if err != nil {
		return nil, err
	}

	response := &SensorResponse{
		ID:           sensor.ID,
		ParkingLotID: sensor.ParkingLotID,
		Status:       sensor.Status,
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
			ID:           sensor.ID,
			ParkingLotID: sensor.ParkingLotID,
			Status:       sensor.Status,
		})
	}

	return response, nil
}
