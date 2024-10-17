package handler

import (
	"net/http"
	"strconv"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/CamiloLeonP/parking-radar/internal/hub"
	"github.com/gin-gonic/gin"
)

// SensorHandler manages sensor operations and WebSocket notifications
type SensorHandler struct {
	SensorUseCase usecase.ISensorUseCase
	WebSocketHub  *hub.WebSocketHub
}

// NewSensorHandler creates a new instance of SensorHandler
func NewSensorHandler(sensorUseCase usecase.ISensorUseCase, wsHub *hub.WebSocketHub) *SensorHandler {
	return &SensorHandler{
		SensorUseCase: sensorUseCase,
		WebSocketHub:  wsHub,
	}
}

// CreateSensor creates a new sensor and notifies clients
func (h *SensorHandler) CreateSensor(c *gin.Context) {
	var req usecase.CreateSensorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.SensorUseCase.CreateSensor(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notify clients about the new sensor creation
	h.NotifyChange("sensor-created", gin.H{
		"device_identifier": req.DeviceIdentifier,
		"sensor_number":     req.SensorNumber,
	})

	c.JSON(http.StatusCreated, gin.H{"status": "sensor created"})
}

// GetSensor retrieves a specific sensor by ID
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

// UpdateSensor updates a sensor and notifies clients
func (h *SensorHandler) UpdateSensor(c *gin.Context) {
	var req usecase.UpdateSensorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensor, err := h.SensorUseCase.GetSensorByDeviceAndNumber(req.DeviceIdentifier, req.SensorNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sensor not found"})
		return
	}

	if err := h.SensorUseCase.UpdateSensor(sensor.ID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notify clients about the sensor update
	h.NotifyChange("sensor-updated", gin.H{
		"id":                sensor.ID,
		"device_identifier": req.DeviceIdentifier,
		"status":            req.Status,
	})

	c.JSON(http.StatusOK, gin.H{"status": "sensor updated"})
}

// DeleteSensor deletes a sensor and notifies clients
func (h *SensorHandler) DeleteSensor(c *gin.Context) {
	sensorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.SensorUseCase.DeleteSensor(uint(sensorID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notify clients about the sensor deletion
	h.NotifyChange("sensor-deleted", gin.H{
		"id": sensorID,
	})

	c.JSON(http.StatusOK, gin.H{"status": "sensor deleted"})
}

// ListSensors lists all sensors for a specific parking lot
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

// NotifyChange sends a unified notification about sensor-related changes
func (h *SensorHandler) NotifyChange(event string, details gin.H) {
	h.WebSocketHub.Broadcast(gin.H{
		"type": "new-change-in-parking",
		"payload": gin.H{
			"event":   event,
			"details": details,
		},
	})
}
