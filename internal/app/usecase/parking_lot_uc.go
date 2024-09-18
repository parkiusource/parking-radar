package usecase

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
)

type IParkingLotUseCase interface {
	CreateParkingLot(req CreateParkingLotRequest) error
	GetParkingLot(parkingLotID uint) (*ParkingLotResponse, error)
	UpdateParkingLot(parkingLotID uint, req UpdateParkingLotRequest) error
	DeleteParkingLot(parkingLotID uint) error
	ListParkingLots() ([]ParkingLotResponse, error)
}

type ParkingLotUseCase struct {
	ParkingLotRepository repository.IParkingLotRepository
	SensorRepository     repository.ISensorRepository
}

type ParkingLotResponse struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Address         string  `json:"address"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	AvailableSpaces uint    `json:"available_spaces"`
}

type CreateParkingLotRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UpdateParkingLotRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewParkingLotUseCase(parkingLotRepo repository.IParkingLotRepository) IParkingLotUseCase {
	return &ParkingLotUseCase{
		ParkingLotRepository: parkingLotRepo,
	}
}

func (uc *ParkingLotUseCase) CreateParkingLot(req CreateParkingLotRequest) error {
	parkingLot := domain.ParkingLot{
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	return uc.ParkingLotRepository.Create(&parkingLot)
}

func (uc *ParkingLotUseCase) GetParkingLot(parkingLotID uint) (*ParkingLotResponse, error) {
	parkingLot, err := uc.ParkingLotRepository.GetByID(parkingLotID)
	if err != nil {
		return nil, err
	}

	sensors, err := uc.SensorRepository.ListByParkingLot(parkingLotID)
	if err != nil {
		return nil, err
	}

	var availableSpaces uint
	for _, sensor := range sensors {
		if sensor.Status == "free" {
			availableSpaces++
		}
	}

	response := &ParkingLotResponse{
		ID:              parkingLot.ID,
		Name:            parkingLot.Name,
		Address:         parkingLot.Address,
		Latitude:        parkingLot.Latitude,
		Longitude:       parkingLot.Longitude,
		AvailableSpaces: availableSpaces,
	}

	return response, nil
}

func (uc *ParkingLotUseCase) UpdateParkingLot(parkingLotID uint, req UpdateParkingLotRequest) error {
	parkingLot, err := uc.ParkingLotRepository.GetByID(parkingLotID)
	if err != nil {
		return err
	}

	parkingLot.Name = req.Name
	parkingLot.Address = req.Address
	parkingLot.Latitude = req.Latitude
	parkingLot.Longitude = req.Longitude

	return uc.ParkingLotRepository.Update(parkingLot)
}

func (uc *ParkingLotUseCase) DeleteParkingLot(parkingLotID uint) error {
	return uc.ParkingLotRepository.Delete(parkingLotID)
}

func (uc *ParkingLotUseCase) ListParkingLots() ([]ParkingLotResponse, error) {
	parkingLots, err := uc.ParkingLotRepository.List()
	if err != nil {
		return nil, err
	}

	var response []ParkingLotResponse
	for _, lot := range parkingLots {
		sensors, err := uc.SensorRepository.ListByParkingLot(lot.ID)
		if err != nil {
			return nil, err
		}
		var availableSpaces uint
		for _, sensor := range sensors {
			if sensor.Status == "free" {
				availableSpaces++
			}
		}

		response = append(response, ParkingLotResponse{
			ID:              lot.ID,
			Name:            lot.Name,
			Address:         lot.Address,
			Latitude:        lot.Latitude,
			Longitude:       lot.Longitude,
			AvailableSpaces: availableSpaces,
		})
	}

	return response, nil
}
