package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/router"
	"github.com/CamiloLeonP/parking-radar/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {

	db.ConnectDatabase()

	err := db.DB.AutoMigrate(&domain.User{}, &domain.ParkingLot{}, &domain.Sensor{}, &domain.Esp32Device{}, &domain.Admin{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database connected and migrated successfully")

	gin.SetMode(gin.ReleaseMode)
	root := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("none port")
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	root.Run(":" + port)

}
