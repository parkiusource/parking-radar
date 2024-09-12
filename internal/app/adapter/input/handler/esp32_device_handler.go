package handler

import (
	"net/http"
	"strconv"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/gin-gonic/gin"
)

const (
	invalidIDError = "invalid id"
)

type Esp32DeviceHandler struct {
	Esp32DeviceUseCase usecase.IEsp32DeviceUseCase
}

func NewEsp32DeviceHandler(esp32DeviceUseCase usecase.IEsp32DeviceUseCase) *Esp32DeviceHandler {
	return &Esp32DeviceHandler{Esp32DeviceUseCase: esp32DeviceUseCase}
}

func (h *Esp32DeviceHandler) CreateEsp32Device(c *gin.Context) {
	var req usecase.CreateEsp32DeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Esp32DeviceUseCase.CreateEsp32Device(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "esp32 device created"})
}

func (h *Esp32DeviceHandler) GetEsp32Device(c *gin.Context) {
	esp32DeviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidIDError})
		return
	}

	esp32Device, err := h.Esp32DeviceUseCase.GetEsp32Device(esp32DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, esp32Device)
}

func (h *Esp32DeviceHandler) UpdateEsp32Device(c *gin.Context) {
	var req usecase.UpdateEsp32DeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	esp32DeviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidIDError})
		return
	}

	err = h.Esp32DeviceUseCase.UpdateEsp32Device(esp32DeviceID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "esp32 device updated"})
}

func (h *Esp32DeviceHandler) DeleteEsp32Device(c *gin.Context) {
	esp32DeviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidIDError})
		return
	}

	err = h.Esp32DeviceUseCase.DeleteEsp32Device(esp32DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "esp32 device deleted"})
}

func (h *Esp32DeviceHandler) ListEsp32Devices(c *gin.Context) {
	esp32Devices, err := h.Esp32DeviceUseCase.ListEsp32Devices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, esp32Devices)
}
