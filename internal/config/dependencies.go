package config

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/output/db"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	db2 "github.com/CamiloLeonP/parking-radar/internal/db"
	"github.com/CamiloLeonP/parking-radar/internal/hub"
	"log"
)

// Handlers stores all the handlers used in the application
type Handlers struct {
	UserHandler        *handler.UserHandler
	ParkingLotHandler  *handler.ParkingLotHandler
	SensorHandler      *handler.SensorHandler
	Esp32DeviceHandler *handler.Esp32DeviceHandler
	WebSocketHandler   *handler.WebSocketHandler
}

// SetupDependencies initializes all dependencies and returns the handlers
func SetupDependencies() *Handlers {
	wsHub := setupWebSocketHub() // Initialize WebSocket hub

	return &Handlers{
		UserHandler:        setupUserHandler(),
		ParkingLotHandler:  setupParkingLotHandler(wsHub),
		SensorHandler:      setupSensorHandler(wsHub),
		Esp32DeviceHandler: setupEsp32DeviceHandler(),
		WebSocketHandler:   setupWebSocketHandler(wsHub),
	}
}

// setupWebSocketHub initializes the WebSocket hub
func setupWebSocketHub() *hub.WebSocketHub {
	h := hub.NewWebSocketHub()
	go func() {
		log.Println("Starting WebSocket hub...")
		h.Run()
	}()
	return h
}

// setupWebSocketHandler initializes the WebSocket handler with the hub
func setupWebSocketHandler(wsHub *hub.WebSocketHub) *handler.WebSocketHandler {
	return handler.NewWebSocketHandler(wsHub)
}

// setupUserHandler initializes the UserHandler
func setupUserHandler() *handler.UserHandler {
	userRepository := &db.UserRepositoryImpl{DB: db2.DB}
	userUseCase := usecase.NewUserUseCase(userRepository)
	return handler.NewUserHandler(userUseCase)
}

// setupParkingLotHandler initializes the ParkingLotHandler with the hub
func setupParkingLotHandler(wsHub *hub.WebSocketHub) *handler.ParkingLotHandler {
	sensorRepository := &db.SensorRepositoryImpl{DB: db2.DB}
	parkingLotRepository := &db.ParkingLotRepositoryImpl{DB: db2.DB}
	parkingLotUseCase := usecase.NewParkingLotUseCase(parkingLotRepository, sensorRepository)
	return handler.NewParkingLotHandler(parkingLotUseCase, wsHub)
}

// setupSensorHandler initializes the SensorHandler with the hub
func setupSensorHandler(wsHub *hub.WebSocketHub) *handler.SensorHandler {
	sensorRepository := &db.SensorRepositoryImpl{DB: db2.DB}
	esp32DeviceRepository := &db.Esp32DeviceRepositoryImpl{DB: db2.DB}
	sensorUseCase := usecase.NewSensorUseCase(sensorRepository, esp32DeviceRepository)
	return handler.NewSensorHandler(sensorUseCase, wsHub)
}

// setupEsp32DeviceHandler initializes the Esp32DeviceHandler
func setupEsp32DeviceHandler() *handler.Esp32DeviceHandler {
	esp32DeviceRepository := &db.Esp32DeviceRepositoryImpl{DB: db2.DB}
	sensorRepository := &db.SensorRepositoryImpl{DB: db2.DB}
	esp32DeviceUseCase := usecase.NewEsp32DeviceUseCase(esp32DeviceRepository, sensorRepository)
	return handler.NewEsp32DeviceHandler(esp32DeviceUseCase)
}
