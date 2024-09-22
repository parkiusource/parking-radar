package config

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/output/db"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	db2 "github.com/CamiloLeonP/parking-radar/internal/db"
)

type Handlers struct {
	UserHandler        *handler.UserHandler
	ParkingLotHandler  *handler.ParkingLotHandler
	SensorHandler      *handler.SensorHandler
	Esp32DeviceHandler *handler.Esp32DeviceHandler
}

func SetupDependencies() *Handlers {
	return &Handlers{
		UserHandler:        setupUserHandler(),
		ParkingLotHandler:  setupParkingLotHandler(),
		SensorHandler:      setupSensorHandler(),
		Esp32DeviceHandler: setupEsp32DeviceHandler(),
	}
}

func setupUserHandler() *handler.UserHandler {
	userRepository := &db.UserRepositoryImpl{DB: db2.DB}
	userUseCase := usecase.NewUserUseCase(userRepository)
	return handler.NewUserHandler(userUseCase)
}

func setupParkingLotHandler() *handler.ParkingLotHandler {
	sensorRepository := &db.SensorRepositoryImpl{DB: db2.DB}
	parkingLotRepository := &db.ParkingLotRepositoryImpl{DB: db2.DB}
	parkingLotUseCase := usecase.NewParkingLotUseCase(parkingLotRepository, sensorRepository)
	return handler.NewParkingLotHandler(parkingLotUseCase)
}

func setupSensorHandler() *handler.SensorHandler {
	sensorRepository := &db.SensorRepositoryImpl{DB: db2.DB}
	esp32DeviceRepository := &db.Esp32DeviceRepositoryImpl{DB: db2.DB}
	sensorUseCase := usecase.NewSensorUseCase(sensorRepository, esp32DeviceRepository)
	return handler.NewSensorHandler(sensorUseCase)
}

func setupEsp32DeviceHandler() *handler.Esp32DeviceHandler {
	esp32DeviceRepository := &db.Esp32DeviceRepositoryImpl{DB: db2.DB}
	sensorRepository := &db.SensorRepositoryImpl{DB: db2.DB}
	esp32DeviceUseCase := usecase.NewEsp32DeviceUseCase(esp32DeviceRepository, sensorRepository)
	return handler.NewEsp32DeviceHandler(esp32DeviceUseCase)
}
