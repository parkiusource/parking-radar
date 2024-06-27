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


func InitHandler(c *gin.Context) {
	
	/*var data domain.InitData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}*/

	Uid := c.GetHeader("UID")
	if Uid == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Uid header missing"})
			c.Abort()
			return
	}

	c.JSON(http.StatusOK, gin.H{
		"Auth": "Autorizado con el dispositivo " + Uid,
	})
}
