package usecase

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
	"time"
)

type IEsp32DeviceUseCase interface {
	CreateEsp32Device(req CreateEsp32DeviceRequest) error
	GetEsp32Device(id uint64) (*Esp32DeviceResponse, error)
	GetEsp32DeviceByIdentifier(identifier string) (*domain.Esp32Device, error)
	UpdateEsp32Device(id uint64, req UpdateEsp32DeviceRequest) error
	DeleteEsp32Device(id uint64) error
	ListEsp32Devices() ([]domain.Esp32Device, error)
}

type Esp32DeviceUseCase struct {
	Esp32DeviceRepository repository.IEsp32DeviceRepository
	SensorRepository      repository.ISensorRepository
}

type CreateEsp32DeviceRequest struct {
	DeviceIdentifier string `json:"device_identifier"`
}

type UpdateEsp32DeviceRequest struct {
	DeviceIdentifier string `json:"device_identifier"`
}

type Esp32DeviceResponse struct {
	ID                uint64          `json:"id"`
	DeviceIdentifier  string          `json:"device_identifier"`
	LastCommunication string          `json:"last_communication"`
	Sensors           []domain.Sensor `json:"sensors"`
}

func NewEsp32DeviceUseCase(esp32DeviceRepo repository.IEsp32DeviceRepository, sensorRepo repository.ISensorRepository) IEsp32DeviceUseCase {
	return &Esp32DeviceUseCase{
		Esp32DeviceRepository: esp32DeviceRepo,
		SensorRepository:      sensorRepo,
	}
}

func (uc *Esp32DeviceUseCase) CreateEsp32Device(req CreateEsp32DeviceRequest) error {
	device := domain.Esp32Device{
		DeviceIdentifier: req.DeviceIdentifier,
	}
	return uc.Esp32DeviceRepository.Create(&device)
}

func (uc *Esp32DeviceUseCase) GetEsp32Device(id uint64) (*Esp32DeviceResponse, error) {
	device, err := uc.Esp32DeviceRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get the sensors for the device
	sensors, err := uc.SensorRepository.ListByEsp32DeviceID(id)
	if err != nil {
		return nil, err
	}

	response := &Esp32DeviceResponse{
		ID:                device.ID,
		DeviceIdentifier:  device.DeviceIdentifier,
		LastCommunication: device.LastCommunication.Format(time.RFC3339), // Formato ISO 8601
		Sensors:           sensors,
	}

	return response, nil
}

func (uc *Esp32DeviceUseCase) GetEsp32DeviceByIdentifier(identifier string) (*domain.Esp32Device, error) {
	return uc.Esp32DeviceRepository.GetByDeviceIdentifier(identifier)
}

func (uc *Esp32DeviceUseCase) UpdateEsp32Device(id uint64, req UpdateEsp32DeviceRequest) error {
	device, err := uc.Esp32DeviceRepository.GetByID(id)
	if err != nil {
		return err
	}

	device.DeviceIdentifier = req.DeviceIdentifier

	return uc.Esp32DeviceRepository.Update(device)
}

func (uc *Esp32DeviceUseCase) DeleteEsp32Device(id uint64) error {
	return uc.Esp32DeviceRepository.Delete(id)
}

func (uc *Esp32DeviceUseCase) ListEsp32Devices() ([]domain.Esp32Device, error) {
	return uc.Esp32DeviceRepository.ListAll()
}
