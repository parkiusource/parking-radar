package handler

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/CamiloLeonP/parking-radar/internal/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

	adminID, _ := helpers.ExtractAdminIDAndRole(c)

	if strings.EqualFold(adminID, "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing admin ID"})
		return
	}

	err := h.AdminUseCase.RegisterAdmin(adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error registering admin"})
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
