package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	middlewares "github.com/CamiloLeonP/parking-radar/internal/middleware"
	"github.com/CamiloLeonP/parking-radar/internal/test/parking/mockgen"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupAdminHandler(mockUseCase *mockgen.MockIAdminUseCase) *gin.Engine {
	adminHandler := NewAdminHandler(mockUseCase)

	r := gin.Default()
	admin := r.Group("/admin")
	admin.Use(middlewares.AuthMiddleware("admin_local", "admin_global"))
	{
		admin.POST("/register", adminHandler.RegisterAdmin)
		admin.PUT("/complete-profile", adminHandler.CompleteAdminProfile)
		admin.GET("/profile", adminHandler.GetAdminProfile)
		admin.GET("/parking-lots", adminHandler.GetParkingLotsByAdmin)
	}

	return r
}

func TestRegisterAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	type RegisterAdminPayload struct {
		User string `json:"user_id"`
	}

	payload := RegisterAdminPayload{
		User: "auth0|6721363081b8547d3f95a976",
	}

	mockUseCase.EXPECT().RegisterAdmin("auth0|6721363081b8547d3f95a976").Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/admin/register", bytes.NewBuffer(body))
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "admin registered successfully", response["status"])
}

func TestRegisterAdmin_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	type RegisterAdminPayload struct {
		User string `json:"user_id"`
	}

	payload := RegisterAdminPayload{
		User: "auth0|6721363081b8547d3f95a976",
	}

	mockUseCase.EXPECT().RegisterAdmin("auth0|6721363081b8547d3f95a976").Return(errors.New("registration failed"))

	w := httptest.NewRecorder()
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/admin/register", bytes.NewBuffer(body))
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "registration failed", response["error"])
}

func TestCompleteAdminProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	profileData := domain.AdminProfileData{
		NIT:          "123456789",
		PhotoURL:     "http://example.com/photo.jpg",
		ContactPhone: "1234567890",
	}

	mockUseCase.EXPECT().CompleteAdminProfile("auth0|6721363081b8547d3f95a976", profileData).Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(profileData)
	req, _ := http.NewRequest("PUT", "/admin/complete-profile", bytes.NewBuffer(body))
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "profile completed", response["status"])
}

func TestCompleteAdminProfile_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	profileData := domain.AdminProfileData{
		NIT:          "123456789",
		PhotoURL:     "http://example.com/photo.jpg",
		ContactPhone: "1234567890",
	}

	mockUseCase.EXPECT().CompleteAdminProfile("auth0|6721363081b8547d3f95a976", profileData).Return(errors.New("failed to complete profile"))

	w := httptest.NewRecorder()
	body, _ := json.Marshal(profileData)
	req, _ := http.NewRequest("PUT", "/admin/complete-profile", bytes.NewBuffer(body))
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "failed to complete profile", response["error"])
}

func TestGetAdminProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	expectedProfile := &domain.Admin{
		Auth0UUID:    "auth0|6721363081b8547d3f95a976",
		NIT:          "123456789",
		PhotoURL:     "http://example.com/photo.jpg",
		ContactPhone: "1234567890",
	}

	mockUseCase.EXPECT().GetAdminProfile("auth0|6721363081b8547d3f95a976").Return(expectedProfile, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/profile", nil)
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Contains(t, response, "profile")
}

func TestGetAdminProfile_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	mockUseCase.EXPECT().GetAdminProfile("auth0|6721363081b8547d3f95a976").Return(nil, errors.New("profile not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/profile", nil)
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "failed to get profile", response["error"])
}

func TestGetParkingLotsByAdmin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	mockParkingLots := []domain.ParkingLot{
		{ID: 1, Name: "Parking Lot 1", Address: "123 Main St"},
		{ID: 2, Name: "Parking Lot 2", Address: "456 Elm St"},
	}
	adminID := "auth0|6721363081b8547d3f95a976"
	mockUseCase.EXPECT().GetParkingLotsByAdmin(adminID).Return(mockParkingLots, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/parking-lots", nil)
	tokenString, _ := generateTestJWT() // Asegúrate de que `generateTestJWT` devuelva un token válido
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)

	parkingLotsData, ok := response["parking_lots"].([]interface{})
	assert.True(t, ok, "la respuesta debe contener una lista de parking_lots")

	for i, lotData := range parkingLotsData {
		lotMap, ok := lotData.(map[string]interface{})
		assert.True(t, ok, "cada parking lot debe ser un mapa")

		assert.Equal(t, mockParkingLots[i].ID, uint(lotMap["id"].(float64)))
		assert.Equal(t, mockParkingLots[i].Name, lotMap["name"])
		assert.Equal(t, mockParkingLots[i].Address, lotMap["address"])
	}
}

func TestGetParkingLotsByAdmin_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIAdminUseCase(ctrl)
	r := setupAdminHandler(mockUseCase)

	adminID := "auth0|6721363081b8547d3f95a976"
	mockUseCase.EXPECT().GetParkingLotsByAdmin(adminID).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/parking-lots", nil)
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return
	}
	assert.Equal(t, "failed to get parking lots", response["error"])
}
