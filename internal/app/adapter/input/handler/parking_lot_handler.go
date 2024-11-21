package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/CamiloLeonP/parking-radar/internal/helpers"
	"github.com/CamiloLeonP/parking-radar/internal/hub"
	"github.com/gin-gonic/gin"
)

const (
	dontHaveAccessToParkingLot = "you don't have access to this parking lot"
	invalidParkingLotID        = "invalid parking lot id"
	invalidRequestBody         = "invalid request body"
)

type ParkingLotHandler struct {
	useCase      usecase.IParkingLotUseCase
	webSocketHub *hub.WebSocketHub
}

func NewParkingLotHandler(useCase usecase.IParkingLotUseCase, wsHub *hub.WebSocketHub) *ParkingLotHandler {
	return &ParkingLotHandler{
		useCase:      useCase,
		webSocketHub: wsHub,
	}
}

func (h *ParkingLotHandler) CreateParkingLot(c *gin.Context) {
	var req usecase.CreateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestBody})
		return
	}

	adminUUID, _ := helpers.ExtractAdminIDAndRole(c)
	if adminUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role for admin"})
		return
	}

	req.AdminUUID = adminUUID
	h.processCreateParkingLot(c, req)
}

func (h *ParkingLotHandler) processCreateParkingLot(c *gin.Context, req usecase.CreateParkingLotRequest) {
	parkingLot, err := h.useCase.CreateParkingLot(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parking lot already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating parking lot"})
		return
	}

	h.notifyChange("parking-lot-created", gin.H{
		"id":      parkingLot.ID,
		"name":    parkingLot.Name,
		"address": parkingLot.Address,
	})

	c.JSON(http.StatusCreated, gin.H{"status": "parking lot created", "id": parkingLot.ID})
}

func (h *ParkingLotHandler) validateAccess(c *gin.Context) (uint, bool) {
	parkingLotID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidParkingLotID})
		return 0, false
	}

	adminUUID, isGlobalAdmin := helpers.ExtractAdminIDAndRole(c)
	if !isGlobalAdmin {
		if _, err := h.useCase.GetParkingLotWithOwnership(uint(parkingLotID), adminUUID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": dontHaveAccessToParkingLot})
			return 0, false
		}
	}

	return uint(parkingLotID), true
}

func (h *ParkingLotHandler) GetParkingLot(c *gin.Context) {
	parkingLotID, ok := h.validateAccess(c)
	if !ok {
		return
	}

	parkingLot, err := h.useCase.GetParkingLot(parkingLotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve parking lot"})
		return
	}

	c.JSON(http.StatusOK, parkingLot)
}

func (h *ParkingLotHandler) UpdateParkingLot(c *gin.Context) {
	parkingLotID, ok := h.validateAccess(c)
	if !ok {
		return
	}

	var req usecase.UpdateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestBody})
		return
	}

	adminUUID, _ := helpers.ExtractAdminIDAndRole(c)
	if err := h.useCase.UpdateParkingLot(parkingLotID, req, adminUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update parking lot"})
		return
	}

	h.notifyChange("parking-lot-updated", gin.H{
		"id":      parkingLotID,
		"name":    req.Name,
		"address": req.Address,
	})

	c.JSON(http.StatusOK, gin.H{"status": "parking lot updated"})
}

func (h *ParkingLotHandler) DeleteParkingLot(c *gin.Context) {
	parkingLotID, ok := h.validateAccess(c)
	if !ok {
		return
	}

	adminUUID, _ := helpers.ExtractAdminIDAndRole(c)
	if err := h.useCase.DeleteParkingLot(parkingLotID, adminUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete parking lot"})
		return
	}

	h.notifyChange("parking-lot-deleted", gin.H{"id": parkingLotID})
	c.JSON(http.StatusOK, gin.H{"status": "parking lot deleted"})
}

func (h *ParkingLotHandler) ListParkingLots(c *gin.Context) {
	parkingLots, err := h.useCase.ListParkingLots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve parking lots"})
		return
	}
	c.JSON(http.StatusOK, parkingLots)
}

func (h *ParkingLotHandler) notifyChange(event string, details gin.H) {
	h.webSocketHub.Broadcast(gin.H{
		"type": "new-change-in-parking",
		"payload": gin.H{
			"event":   event,
			"details": details,
		},
	})
}
