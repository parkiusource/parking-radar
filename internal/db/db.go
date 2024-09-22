package db

import (
	"log"

	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	// Conectar a la base de datos SQLite (se creará el archivo database.db si no existe)
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

}
