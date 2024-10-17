package router

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/config"
	middlewares "github.com/CamiloLeonP/parking-radar/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "https://parking-radar-frontend.vercel.app"
	}

	r.Use(middlewares.CORSMiddleware(allowedOrigin))

	handlers := config.SetupDependencies()

	// Rutas para usuarios
	users := r.Group("/users")
	{
		users.POST("/register", handlers.UserHandler.Register)
		users.GET("/:id", handlers.UserHandler.GetUserByID)
		users.PUT("/:id", handlers.UserHandler.UpdateUser)
		users.DELETE("/:id", handlers.UserHandler.DeleteUser)
	}

	// Rutas para parking lots
	parkingLots := r.Group("/parking-lots")
	{
		parkingLots.POST("/", handlers.ParkingLotHandler.CreateParkingLot)
		parkingLots.GET("/", handlers.ParkingLotHandler.ListParkingLots)
		parkingLots.GET("/:id", handlers.ParkingLotHandler.GetParkingLot)
		parkingLots.PUT("/:id", handlers.ParkingLotHandler.UpdateParkingLot)
		parkingLots.DELETE("/:id", handlers.ParkingLotHandler.DeleteParkingLot)
	}

	// Rutas para sensores
	sensors := r.Group("/sensors")
	{
		sensors.POST("/", handlers.SensorHandler.CreateSensor)
		sensors.GET("/", handlers.SensorHandler.ListSensors)
		sensors.GET("/:id", handlers.SensorHandler.GetSensor)
		sensors.PUT("/:sensor_number", handlers.SensorHandler.UpdateSensor)
		sensors.DELETE("/:id", handlers.SensorHandler.DeleteSensor)
	}

	// Rutas para ESP32 devices
	esp32Devices := r.Group("/esp32-devices")
	{
		esp32Devices.POST("/register", handlers.Esp32DeviceHandler.CreateEsp32Device)
		esp32Devices.GET("/list", handlers.Esp32DeviceHandler.ListEsp32Devices)
		esp32Devices.GET("/:id", handlers.Esp32DeviceHandler.GetEsp32Device)
		esp32Devices.PUT("/:id", handlers.Esp32DeviceHandler.UpdateEsp32Device)
		esp32Devices.DELETE("/:id", handlers.Esp32DeviceHandler.DeleteEsp32Device)
	}

	r.POST("/ping", handler.PinHandler)

	r.GET("/init", handler.AuthMiddleware(), handler.InitHandler)

	return r
}
