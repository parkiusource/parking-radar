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
	}

	r.POST("/ping", handler.PinHandler)

	r.GET("/init", handler.AuthMiddleware(), handler.InitHandler)

	return r
}
