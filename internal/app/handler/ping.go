package handler

import (
	"net/http"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/gin-gonic/gin"
)

func PinHandler(c *gin.Context) {
	var data domain.SensorData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "received data",
		"yourData": data.Test,
	})
}
