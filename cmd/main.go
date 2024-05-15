package main

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/ping", handler.PinHandler)

	router.Run(":8080")

}
