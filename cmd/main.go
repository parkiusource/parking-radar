package main

import (
	"os"

	"github.com/CamiloLeonP/parking-radar/internal/app/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	router.POST("/ping", handler.PinHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	router.Run(":" + port)

}
