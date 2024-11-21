package handler

import (
	"net/http"
	"strings"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/CamiloLeonP/parking-radar/internal/helpers"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminUseCase usecase.IAdminUseCase
}

func NewAdminHandler(adminUseCase usecase.IAdminUseCase) *AdminHandler {
	return &AdminHandler{
		AdminUseCase: adminUseCase,
	}
}

func (h *AdminHandler) RegisterAdmin(c *gin.Context) {
	// Define a struct to parse the request body
	type RegisterAdminPayload struct {
		User string `json:"user_id"`
	}

	var payload RegisterAdminPayload

	// Bind JSON payload to the struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Extract adminID from the payload
	adminID := payload.User

	if strings.EqualFold(adminID, "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing admin ID"})
		return
	}

	// Call the use case to register the admin
	err := h.AdminUseCase.RegisterAdmin(adminID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "admin registered successfully"})
}

func (h *AdminHandler) CompleteAdminProfile(c *gin.Context) {
	adminID, _ := helpers.ExtractAdminIDAndRole(c)

	var profileData domain.AdminProfileData
	if err := c.ShouldBindJSON(&profileData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AdminUseCase.CompleteAdminProfile(adminID, profileData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "profile completed"})
}

func (h *AdminHandler) GetAdminProfile(c *gin.Context) {
	adminID, _ := helpers.ExtractAdminIDAndRole(c)

	profile, err := h.AdminUseCase.GetAdminProfile(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (h *AdminHandler) GetParkingLotsByAdmin(c *gin.Context) {
	adminID, _ := helpers.ExtractAdminIDAndRole(c)

	parkingLots, err := h.AdminUseCase.GetParkingLotsByAdmin(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get parking lots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"parking_lots": parkingLots})
}
