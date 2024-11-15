package handler

import (
	"github.com/CamiloLeonP/parking-radar/internal/helpers"
	"net/http"
	"strconv"
	"strings"

	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/CamiloLeonP/parking-radar/internal/hub"
	"github.com/gin-gonic/gin"
)

const dontHaveAccessToParkingLot = "you don't have access to this parking lot"

// ParkingLotHandler manages parking lot operations and WebSocket notifications.
type ParkingLotHandler struct {
	ParkingLotUseCase usecase.IParkingLotUseCase
	WebSocketHub      *hub.WebSocketHub // WebSocket hub dependency.
}

// NewParkingLotHandler creates a new instance of ParkingLotHandler.
func NewParkingLotHandler(ParkingLotUseCase usecase.IParkingLotUseCase, wsHub *hub.WebSocketHub) *ParkingLotHandler {
	return &ParkingLotHandler{
		ParkingLotUseCase: ParkingLotUseCase,
		WebSocketHub:      wsHub,
	}
}

// CreateParkingLot creates a new parking lot and notifies clients.
func (h *ParkingLotHandler) CreateParkingLot(c *gin.Context) {
	var req usecase.CreateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	AdminUUID, _ := helpers.ExtractAdminIDAndRole(c)

	if strings.EqualFold(AdminUUID, "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role for admin"})
		return
	}

	req.AdminUUID = AdminUUID

	h.processCreateParkingLot(c, req)
}

func (h *ParkingLotHandler) processCreateParkingLot(c *gin.Context, req usecase.CreateParkingLotRequest) {
	parkingLot, err := h.ParkingLotUseCase.CreateParkingLot(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "parking lot already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating parking lot"})
		return
	}

	h.NotifyChange("parking-lot-created", gin.H{
		"id":      parkingLot.ID,
		"name":    parkingLot.Name,
		"address": parkingLot.Address,
	})

	c.JSON(http.StatusCreated, gin.H{"status": "parking lot created", "id": parkingLot.ID})
}

// GetParkingLot retrieves a specific parking lot by ID.
func (h *ParkingLotHandler) GetParkingLot(c *gin.Context) {
	adminUUID, isGlobalAdmin := helpers.ExtractAdminIDAndRole(c)

	parkingLotID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parking lot id"})
		return
	}

	if !isGlobalAdmin {
		if _, err := h.ParkingLotUseCase.GetParkingLotWithOwnership(uint(parkingLotID), adminUUID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": dontHaveAccessToParkingLot})
			return
		}
	}

	parkingLot, err := h.ParkingLotUseCase.GetParkingLot(uint(parkingLotID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve parking lot"})
		return
	}

	c.JSON(http.StatusOK, parkingLot)
}

// UpdateParkingLot updates a parking lot and notifies clients.
func (h *ParkingLotHandler) UpdateParkingLot(c *gin.Context) {
	adminUUID, isGlobalAdmin := helpers.ExtractAdminIDAndRole(c)

	parkingLotID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parking lot id"})
		return
	}

	if !isGlobalAdmin {
		if _, err := h.ParkingLotUseCase.GetParkingLotWithOwnership(uint(parkingLotID), adminUUID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": dontHaveAccessToParkingLot})
			return
		}
	}

	var req usecase.UpdateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.ParkingLotUseCase.UpdateParkingLot(uint(parkingLotID), req, adminUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update parking lot"})
		return
	}

	h.NotifyChange("parking-lot-updated", gin.H{
		"id":      parkingLotID,
		"name":    req.Name,
		"address": req.Address,
	})

	c.JSON(http.StatusOK, gin.H{"status": "parking lot updated"})
}

// DeleteParkingLot deletes a parking lot and notifies clients.
func (h *ParkingLotHandler) DeleteParkingLot(c *gin.Context) {
	adminUUID, isGlobalAdmin := helpers.ExtractAdminIDAndRole(c)

	parkingLotID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parking lot id"})
		return
	}

	if !isGlobalAdmin {
		if _, err := h.ParkingLotUseCase.GetParkingLotWithOwnership(uint(parkingLotID), adminUUID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": dontHaveAccessToParkingLot})
			return
		}
	}

	if err := h.ParkingLotUseCase.DeleteParkingLot(uint(parkingLotID), adminUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete parking lot"})
		return
	}

	h.NotifyChange("parking-lot-deleted", gin.H{
		"id": parkingLotID,
	})

	c.JSON(http.StatusOK, gin.H{"status": "parking lot deleted"})
}

// ListParkingLots retrieves all parking lots.
func (h *ParkingLotHandler) ListParkingLots(c *gin.Context) {
	parkingLots, err := h.ParkingLotUseCase.ListParkingLots()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve parking lots"})
		return
	}
	c.JSON(http.StatusOK, parkingLots)
}

// NotifyChange sends a unified notification to all connected clients.
func (h *ParkingLotHandler) NotifyChange(event string, details gin.H) {
	h.WebSocketHub.Broadcast(gin.H{
		"type": "new-change-in-parking",
		"payload": gin.H{
			"event":   event,
			"details": details,
		},
	})
}
