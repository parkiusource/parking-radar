package main

import (
	"database/sql"
	"log"
	"os"
	"github.com/CamiloLeonP/parking-radar/internal/app/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var dbName = "smartparkinguan"
var dbUser = "admin"
var dbPassword = "admin123"
var dbEndpoint = "smartparkinguan.cdywy268kwq7.us-east-2.rds.amazonaws.com" //HOST
var dbEndpointPort = "3306" // 3306 is Default port if not specified

func main() {

	initDB()
	defer db.Close()
	
	gin.SetMode(gin.ReleaseMode)
	root := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("none port")
		port = "8080" // 8080 is Default port if not specified
	}

	log.Printf("Server running on port %s", port)
	root.Run(":" + port)

}


func initDB() {
		var err error
		dsn := dbUser + ":"+ dbPassword + "@tcp(" + dbEndpoint + ":" + dbEndpointPort + ")/" + dbName + "?parseTime=true"
		log.Printf(dsn)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error opening database: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		log.Printf("Conexon exitosa")
}
