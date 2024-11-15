package usecase

import (
	"errors"
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
)

//go:generate mockgen -source=./parking_lot_uc.go -destination=./../../test/parking/mocks/mock_parking_lot_uc.go -package=mockgen
type IParkingLotUseCase interface {
	CreateParkingLot(req CreateParkingLotRequest) (*ParkingLotResponse, error)
	GetParkingLot(parkingLotID uint) (*ParkingLotResponse, error)
	GetParkingLotWithOwnership(parkingLotID uint, adminUUID string) (*ParkingLotResponse, error)
	UpdateParkingLot(parkingLotID uint, req UpdateParkingLotRequest, adminUUID string) error
	DeleteParkingLot(parkingLotID uint, adminUUID string) error
	ListParkingLots() ([]ParkingLotResponse, error)
}

type ParkingLotUseCase struct {
	ParkingLotRepository repository.IParkingLotRepository
	SensorRepository     repository.ISensorRepository
	AdminRepository      repository.IAdminRepository
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
	AdminUUID string  `json:"admin_uuid"`
}

type UpdateParkingLotRequest struct {
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// NewParkingLotUseCase creates a new instance of ParkingLotUseCase.
func NewParkingLotUseCase(parkingLotRepo repository.IParkingLotRepository, sensorRepository repository.ISensorRepository, adminRepository repository.IAdminRepository) IParkingLotUseCase {
	return &ParkingLotUseCase{
		ParkingLotRepository: parkingLotRepo,
		SensorRepository:     sensorRepository,
		AdminRepository:      adminRepository,
	}
}

// CreateParkingLot creates a new parking lot.
func (uc *ParkingLotUseCase) CreateParkingLot(req CreateParkingLotRequest) (*ParkingLotResponse, error) {

	admin, err := uc.AdminRepository.FindByAuth0UUID(req.AdminUUID)
	if err != nil {
		return nil, err
	}

	parkingLot := domain.ParkingLot{
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		AdminID:   admin.ID,
	}

	if err := uc.ParkingLotRepository.Create(&parkingLot); err != nil {
		return nil, err
	}

	return &ParkingLotResponse{
		ID:        parkingLot.ID,
		Name:      parkingLot.Name,
		Address:   parkingLot.Address,
		Latitude:  parkingLot.Latitude,
		Longitude: parkingLot.Longitude,
	}, nil
}

// GetParkingLot retrieves a parking lot by ID.
func (uc *ParkingLotUseCase) GetParkingLot(parkingLotID uint) (*ParkingLotResponse, error) {
	parkingLot, err := uc.ParkingLotRepository.GetByID(parkingLotID)
	if err != nil {
		return nil, err
	}

	sensors, err := uc.SensorRepository.ListByParkingLot(parkingLotID)
	if err != nil {
		return nil, err
	}

	availableSpaces := countAvailableSpaces(sensors)

	return &ParkingLotResponse{
		ID:              parkingLot.ID,
		Name:            parkingLot.Name,
		Address:         parkingLot.Address,
		Latitude:        parkingLot.Latitude,
		Longitude:       parkingLot.Longitude,
		AvailableSpaces: availableSpaces,
	}, nil
}

// GetParkingLotWithOwnership retrieves a parking lot if it belongs to the admin.
func (uc *ParkingLotUseCase) GetParkingLotWithOwnership(parkingLotID uint, adminUUID string) (*ParkingLotResponse, error) {
	admin, err := uc.AdminRepository.FindByAuth0UUID(adminUUID)
	if err != nil {
		return nil, err
	}

	parkingLot, err := uc.ParkingLotRepository.GetByIDWithAdmin(parkingLotID, admin.ID)
	if err != nil {
		return nil, err
	}

	sensors, err := uc.SensorRepository.ListByParkingLot(parkingLotID)
	if err != nil {
		return nil, err
	}

	availableSpaces := countAvailableSpaces(sensors)

	return &ParkingLotResponse{
		ID:              parkingLot.ID,
		Name:            parkingLot.Name,
		Address:         parkingLot.Address,
		Latitude:        parkingLot.Latitude,
		Longitude:       parkingLot.Longitude,
		AvailableSpaces: availableSpaces,
	}, nil
}

// UpdateParkingLot updates a parking lot.
func (uc *ParkingLotUseCase) UpdateParkingLot(parkingLotID uint, req UpdateParkingLotRequest, adminUUID string) error {

	admin, err := uc.AdminRepository.FindByAuth0UUID(adminUUID)
	if err != nil {
		return err
	}

	parkingLot, err := uc.ParkingLotRepository.GetByIDWithAdmin(parkingLotID, admin.ID)
	if err != nil {
		return err
	}

	parkingLot.Name = req.Name
	parkingLot.Address = req.Address
	parkingLot.Latitude = req.Latitude
	parkingLot.Longitude = req.Longitude

	return uc.ParkingLotRepository.Update(parkingLot)
}

// DeleteParkingLot deletes a parking lot with ownership validation.
func (uc *ParkingLotUseCase) DeleteParkingLot(parkingLotID uint, adminUUID string) error {

	admin, err := uc.AdminRepository.FindByAuth0UUID(adminUUID)
	if err != nil {
		return err
	}

	if _, err := uc.ParkingLotRepository.GetByIDWithAdmin(parkingLotID, admin.ID); err != nil {
		return errors.New("forbidden: you don't have access to this parking lot")
	}
	return uc.ParkingLotRepository.Delete(parkingLotID)
}

// ListParkingLots retrieves all parking lots with their available spaces.
func (uc *ParkingLotUseCase) ListParkingLots() ([]ParkingLotResponse, error) {
	parkingLots, err := uc.ParkingLotRepository.List()
	if err != nil {
		return nil, err
	}

	sensorMap, err := uc.SensorRepository.ListGroupedByParkingLot()
	if err != nil {
		return nil, err
	}

	var response []ParkingLotResponse
	for _, lot := range parkingLots {
		availableSpaces := sensorMap[lot.ID]

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

// Helper function to count available spaces.
func countAvailableSpaces(sensors []domain.Sensor) uint {
	var availableSpaces uint
	for _, sensor := range sensors {
		if sensor.Status == "free" {
			availableSpaces++
		}
	}
	return availableSpaces
}
