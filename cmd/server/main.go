package main

import (
	"log"
	"os"
	"github.com/CamiloLeonP/parking-radar/internal/app/router"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	root := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("none port")
		port = "8080" // Default port if not specified
	}

	log.Printf("Server running on port %s", port)
	root.Run(":" + port)

}
