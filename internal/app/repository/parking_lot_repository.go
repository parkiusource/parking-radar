package repository

import "github.com/CamiloLeonP/parking-radar/internal/app/domain"

//go:generate mockgen -source=./parking_lot_repository.go -destination=./../../test/shared/mocks/mock_parking_lot_repository.go -package=mockgen
type IParkingLotRepository interface {
	Create(parkingLot *domain.ParkingLot) error
	GetByID(id uint) (*domain.ParkingLot, error)
	Update(parkingLot *domain.ParkingLot) error
	Delete(id uint) error
	List() ([]domain.ParkingLot, error)
	GetByIDWithAdmin(parkingLotID uint, adminID uint) (*domain.ParkingLot, error)
	FindByAdminID(adminID uint) ([]domain.ParkingLot, error)
}
