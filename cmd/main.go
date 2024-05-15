package main

import (
	"log"
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
		log.Println("none port")
		port = "8080" // Default port if not specified
	}

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)

}
