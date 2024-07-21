package router

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/adapter/input/handler"
	"github.com/CamiloLeonP/parking-radar/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	handlers := config.SetupDependencies()
	users := r.Group("/users")
	{
		users.POST("/register", handlers.UserHandler.Register)
		users.GET("/:id", handlers.UserHandler.GetUserByID)
		users.PUT("/:id", handlers.UserHandler.UpdateUser)
		users.DELETE("/:id", handlers.UserHandler.DeleteUser)
	}

	parkingLots := r.Group("/parkinglots")
	{
		parkingLots.POST("/", handlers.ParkingLotHandler.CreateParkingLot)
		parkingLots.GET("/", handlers.ParkingLotHandler.ListParkingLots)
		parkingLots.GET("/:id", handlers.ParkingLotHandler.GetParkingLot)
		parkingLots.PUT("/:id", handlers.ParkingLotHandler.UpdateParkingLot)
		parkingLots.DELETE("/:id", handlers.ParkingLotHandler.DeleteParkingLot)
	}

	r.POST("/ping", handler.PinHandler)

	r.GET("/init", handler.AuthMiddleware(), handler.InitHandler)

	return r
}
