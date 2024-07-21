package handler

import (
	"net/http"
	"strconv"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/gin-gonic/gin"
)

type ParkingLotHandler struct {
	ParkingLotUseCase usecase.ParkingLotUseCase
}

func NewParkingLotHandler(ParkingLotUseCase usecase.ParkingLotUseCase) *ParkingLotHandler {
	return &ParkingLotHandler{ParkingLotUseCase: ParkingLotUseCase}
}

func (h *ParkingLotHandler) CreateParkingLot(c *gin.Context) {
	var req usecase.CreateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ParkingLotUseCase.CreateParkingLot(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "parking lot created"})
}

func (h *ParkingLotHandler) GetParkingLot(c *gin.Context) {
	parkingLotID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	parkingLot, err := h.ParkingLotUseCase.GetParkingLot(uint(parkingLotID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, parkingLot)
}

func (h *ParkingLotHandler) UpdateParkingLot(c *gin.Context) {
	var req usecase.UpdateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parkingLotID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.ParkingLotUseCase.UpdateParkingLot(uint(parkingLotID), req)
	if err != nil {
		if err.Error() == "no available spaces left to occupy" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "parking lot updated"})
}

func (h *ParkingLotHandler) DeleteParkingLot(c *gin.Context) {
	parkingLotID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.ParkingLotUseCase.DeleteParkingLot(uint(parkingLotID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "parking lot deleted"})
}

func (h *ParkingLotHandler) ListParkingLots(c *gin.Context) {
	parkingLots, err := h.ParkingLotUseCase.ListParkingLots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, parkingLots)
}
