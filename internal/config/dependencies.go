package config

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/output/db"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
)

type Handlers struct {
	UserHandler       *handler.UserHandler
	ParkingLotHandler *handler.ParkingLotHandler
}

func SetupDependencies() *Handlers {
	// Configuración para User
	userRepository := &db.UserRepository{}
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	// Configuración para ParkingLot
	parkingLotRepository := &db.ParkingLotRepository{}
	parkingLotUsecase := &usecase.ParkingLotUseCase{ParkingLotRepository: parkingLotRepository}
	parkingLotHandler := handler.NewParkingLotHandler(*parkingLotUsecase)

	return &Handlers{
		UserHandler:       userHandler,
		ParkingLotHandler: parkingLotHandler,
	}
}
