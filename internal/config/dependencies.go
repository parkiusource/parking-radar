package config

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/output/db"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	db2 "github.com/CamiloLeonP/parking-radar/internal/db"
)

type Handlers struct {
	UserHandler       *handler.UserHandler
	ParkingLotHandler *handler.ParkingLotHandler
}

func SetupDependencies() *Handlers {
	// Configuración para User
	userRepository := &db.UserRepositoryImpl{}
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	// Configuración para ParkingLot
	parkingLotRepository := &db.ParkingLotRepositoryImpl{DB: db2.DB}
	parkingLotUseCase := usecase.NewParkingLotUseCase(parkingLotRepository)
	parkingLotHandler := handler.NewParkingLotHandler(parkingLotUseCase)

	return &Handlers{
		UserHandler:       userHandler,
		ParkingLotHandler: parkingLotHandler,
	}
}
