package usecase

import (
	"errors"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
)

type ParkingLotUseCase struct {
	ParkingLotRepository repository.ParkingLotRepository
}

type CreateParkingLotRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	TotalSpaces int    `json:"total_spaces"`
}

type UpdateParkingLotRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	TotalSpaces int    `json:"total_spaces"`
}

func (uc *ParkingLotUseCase) CreateParkingLot(req CreateParkingLotRequest) error {
	parkingLot := domain.ParkingLot{
		Name:        req.Name,
		Location:    req.Location,
		TotalSpaces: req.TotalSpaces,
	}

	return uc.ParkingLotRepository.Create(&parkingLot)
}

func (uc *ParkingLotUseCase) GetParkingLot(parkingLotID uint) (*domain.ParkingLot, error) {
	return uc.ParkingLotRepository.GetByID(parkingLotID)
}

func (uc *ParkingLotUseCase) UpdateParkingLot(parkingLotID uint, req UpdateParkingLotRequest) error {
	parkingLot, err := uc.ParkingLotRepository.GetByID(parkingLotID)
	if err != nil {
		return err
	}

	parkingLot.Name = req.Name
	parkingLot.Location = req.Location

	// Actualizar availableSpaces solo si totalSpaces cambia
	if parkingLot.TotalSpaces != req.TotalSpaces {
		difference := req.TotalSpaces - parkingLot.TotalSpaces
		parkingLot.AvailableSpaces += difference

		if parkingLot.AvailableSpaces < 0 {
			return errors.New(("no available spaces left to occupy"))
		}
	}

	parkingLot.TotalSpaces = req.TotalSpaces

	return uc.ParkingLotRepository.Update(parkingLot)
}

func (uc *ParkingLotUseCase) DeleteParkingLot(parkingLotID uint) error {
	return uc.ParkingLotRepository.Delete(parkingLotID)
}

func (uc *ParkingLotUseCase) ListParkingLots() ([]domain.ParkingLot, error) {
	return uc.ParkingLotRepository.List()
}
