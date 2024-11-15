package router

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/config"
	middlewares "github.com/CamiloLeonP/parking-radar/internal/middleware"
	"github.com/gin-gonic/gin"
)

const (
	REGISTER = "/register"
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
		users.POST(REGISTER, handlers.UserHandler.Register)
		users.GET("/:id", handlers.UserHandler.GetUserByID)
		users.PUT("/:id", handlers.UserHandler.UpdateUser)
		users.DELETE("/:id", handlers.UserHandler.DeleteUser)
	}

	// Group for parking lots
	publicParkingLots := r.Group("/parking-lots")
	{
		publicParkingLots.GET("/", handlers.ParkingLotHandler.ListParkingLots)
	}

	// Group for protected parking lots
	protectedParkingLots := r.Group("/parking-lots")
	protectedParkingLots.Use(middlewares.AuthMiddleware("admin_local", "admin_global"))
	{
		protectedParkingLots.POST("/", handlers.ParkingLotHandler.CreateParkingLot)
		protectedParkingLots.GET("/:id", handlers.ParkingLotHandler.GetParkingLot)
		protectedParkingLots.PUT("/:id", handlers.ParkingLotHandler.UpdateParkingLot)
		protectedParkingLots.DELETE("/:id", handlers.ParkingLotHandler.DeleteParkingLot)
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

	// Group for register Admin
	admins := r.Group("/admins")
	admins.Use(middlewares.AuthMiddleware("admin_default", "admin_global"))
	{
		admins.POST(REGISTER, handlers.AdminHandler.RegisterAdmin)
	}

	// Group for protected Admin Profile
	protectedAdmins := r.Group("/admins")
	protectedAdmins.Use(middlewares.AuthMiddleware("admin_local", "admin_global"))
	{
		protectedAdmins.GET("/parking-lots", handlers.AdminHandler.GetParkingLotsByAdmin)
		protectedAdmins.POST("/complete-profile", handlers.AdminHandler.CompleteAdminProfile)
		protectedAdmins.GET("/profile", handlers.AdminHandler.GetAdminProfile)

	}

	// Routes for esp32 devices
	esp32Devices := r.Group("/esp32-devices")
	{
		esp32Devices.POST(REGISTER, handlers.Esp32DeviceHandler.CreateEsp32Device)
		esp32Devices.GET("/list", handlers.Esp32DeviceHandler.ListEsp32Devices)
		esp32Devices.GET("/:id", handlers.Esp32DeviceHandler.GetEsp32Device)
		esp32Devices.PUT("/:id", handlers.Esp32DeviceHandler.UpdateEsp32Device)
		esp32Devices.DELETE("/:id", handlers.Esp32DeviceHandler.DeleteEsp32Device)
	}

	r.POST("/ping", handler.PinHandler)

	r.GET("/init", middlewares.AuthMiddleware(), handler.InitHandler)

	return r
}
