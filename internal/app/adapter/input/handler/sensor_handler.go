package handler

import (
	"net/http"
	"strconv"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/gin-gonic/gin"
)

type SensorHandler struct {
	SensorUseCase usecase.ISensorUseCase
}

func NewSensorHandler(sensorUseCase usecase.ISensorUseCase) *SensorHandler {
	return &SensorHandler{SensorUseCase: sensorUseCase}
}

func (h *SensorHandler) CreateSensor(c *gin.Context) {
	var req usecase.CreateSensorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.SensorUseCase.CreateSensor(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "sensor created"})
}

func (h *SensorHandler) GetSensor(c *gin.Context) {
	sensorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sensor, err := h.SensorUseCase.GetSensor(uint(sensorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sensor)
}

func (h *SensorHandler) UpdateSensor(c *gin.Context) {
	var req usecase.UpdateSensorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.SensorUseCase.UpdateSensor(uint(sensorID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sensor updated"})
}

func (h *SensorHandler) DeleteSensor(c *gin.Context) {
	sensorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.SensorUseCase.DeleteSensor(uint(sensorID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "sensor deleted"})
}

func (h *SensorHandler) ListSensors(c *gin.Context) {
	parkingIDStr := c.Query("parkingID")
	parkingID, err := strconv.ParseUint(parkingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parkingID"})
		return
	}

	sensors, err := h.SensorUseCase.ListSensorsByParkingLot(uint(parkingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sensors)
}
