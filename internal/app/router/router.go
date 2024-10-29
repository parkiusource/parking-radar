package router

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/config"
	middlewares "github.com/CamiloLeonP/parking-radar/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = true

	r.Use(middlewares.CORSMiddleware())

	handlers := config.SetupDependencies()

	r.GET("/ws", handlers.WebSocketHandler.HandleConnection)

	// Routes for users
	users := r.Group("/users")
	{
		users.POST("/register", handlers.UserHandler.Register)
		users.GET("/:id", handlers.UserHandler.GetUserByID)
		users.PUT("/:id", handlers.UserHandler.UpdateUser)
		users.DELETE("/:id", handlers.UserHandler.DeleteUser)
	}

	// Routes for parking lots

	parkingLots := r.Group("/parking-lots")
	parkingLots.Use(middlewares.AuthMiddleware("admin_local", "admin_global"))
	r.GET("/parking-lots", handlers.ParkingLotHandler.ListParkingLots)
	{
		parkingLots.POST("/", handlers.ParkingLotHandler.CreateParkingLot)
		parkingLots.PUT("/:id", handlers.ParkingLotHandler.UpdateParkingLot)
		parkingLots.DELETE("/:id", handlers.ParkingLotHandler.DeleteParkingLot)
		parkingLots.GET("/:id", handlers.ParkingLotHandler.GetParkingLot)
	}

	// Routes for sensors
	sensors := r.Group("/sensors")
	{
		sensors.POST("/", handlers.SensorHandler.CreateSensor)
		sensors.GET("/", handlers.SensorHandler.ListSensors)
		sensors.GET("/:id", handlers.SensorHandler.GetSensor)
		sensors.PUT("/:sensor_number", handlers.SensorHandler.UpdateSensor)
		sensors.DELETE("/:id", handlers.SensorHandler.DeleteSensor)
	}

	// Routes for esp32 devices
	esp32Devices := r.Group("/esp32-devices")
	{
		esp32Devices.POST("/register", handlers.Esp32DeviceHandler.CreateEsp32Device)
		esp32Devices.GET("/list", handlers.Esp32DeviceHandler.ListEsp32Devices)
		esp32Devices.GET("/:id", handlers.Esp32DeviceHandler.GetEsp32Device)
		esp32Devices.PUT("/:id", handlers.Esp32DeviceHandler.UpdateEsp32Device)
		esp32Devices.DELETE("/:id", handlers.Esp32DeviceHandler.DeleteEsp32Device)
	}

	r.POST("/ping", handler.PinHandler)

	r.GET("/init", middlewares.AuthMiddleware(), handler.InitHandler)

	return r
}
