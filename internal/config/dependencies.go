package config

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/output/db"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
)

type Handlers struct {
	UserHandler *handler.UserHandler
	// Agrega otros handlers aquí
}

func SetupDependencies() *Handlers {
	// Configuración para User
	userRepository := &db.UserRepository{}
	userUsecase := &usecase.UserUseCase{UserRepository: userRepository}
	userHandler := handler.NewUserHandler(*userUsecase)

	/*
		// Configuración para ParkingLot
		parkingRepository := &db.ParkingRepository{}
		parkingUsecase := &usecase.ParkingUsecase{ParkingRepository: parkingRepository}
		parkingHandler := rest.NewParkingHandler(parkingUsecase)*/

	return &Handlers{
		UserHandler: userHandler,
		// Agrega otros handlers aquí
	}
}
