package handler

import (
	"bytes"
	"encoding/json"
	"github.com/CamiloLeonP/parking-radar/internal/app/usecase"
	"github.com/CamiloLeonP/parking-radar/internal/hub"
	middlewares "github.com/CamiloLeonP/parking-radar/internal/middleware"
	"github.com/CamiloLeonP/parking-radar/internal/test/parking/mockgen"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// setupParkingLotHandler initializes the handler with middleware and mocks.
func setupParkingLotHandler(mockUseCase *mockgen.MockIParkingLotUseCase) (*gin.Engine, *hub.WebSocketHub) {
	wsHub := hub.NewWebSocketHub()
	go wsHub.Run()

	parkingLotHandler := NewParkingLotHandler(mockUseCase, wsHub)

	r := gin.Default()
	parkingLots := r.Group("/parkinglots")
	parkingLots.Use(middlewares.AuthMiddleware("admin_local", "admin_global"))
	{
		parkingLots.POST("/", parkingLotHandler.CreateParkingLot)
		parkingLots.GET("/:id", parkingLotHandler.GetParkingLot)
		parkingLots.PUT("/:id", parkingLotHandler.UpdateParkingLot)
		parkingLots.DELETE("/:id", parkingLotHandler.DeleteParkingLot)
		parkingLots.GET("/", parkingLotHandler.ListParkingLots)
	}

	return r, wsHub
}

// Helper to generate a JWT for testing.
func generateTestJWT() (string, error) {
	privateKeyPEM := os.Getenv("RSA_PRIVATE_KEY")

	return privateKeyPEM, nil
}

// Test for creating a parking lot.
func TestCreateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIParkingLotUseCase(ctrl)
	r, wsHub := setupParkingLotHandler(mockUseCase)
	defer wsHub.Stop()

	responseParking := &usecase.ParkingLotResponse{
		ID: 1, Name: "Test Lot", Address: "123 Test St", Latitude: 12.34, Longitude: 56.78, AvailableSpaces: 10,
	}
	mockUseCase.EXPECT().CreateParkingLot(gomock.Any()).Return(responseParking, nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(usecase.CreateParkingLotRequest{
		Name:      "Test Lot",
		Address:   "123 Test St",
		Latitude:  12.34,
		Longitude: 56.78,
	})
	req, _ := http.NewRequest("POST", "/parkinglots/", bytes.NewBuffer(body))
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "parking lot created", response["status"])
}

func TestGetParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIParkingLotUseCase(ctrl)
	r, wsHub := setupParkingLotHandler(mockUseCase)
	defer wsHub.Stop()

	mockParkingLot := &usecase.ParkingLotResponse{
		ID: 1, Name: "Lot 1", Address: "123 Test St", Latitude: 12.34, Longitude: 56.78, AvailableSpaces: 10,
	}

	mockUseCase.EXPECT().GetParkingLotWithOwnership(uint(1), "auth0|672048d1f0cb0992821786c5").Return(mockParkingLot, nil).Times(0)

	mockUseCase.EXPECT().GetParkingLot(uint(1)).Return(mockParkingLot, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/parkinglots/1", nil)
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseParkingLot usecase.ParkingLotResponse
	err := json.Unmarshal(w.Body.Bytes(), &responseParkingLot)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	assert.Equal(t, mockParkingLot, &responseParkingLot)
}

func TestUpdateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIParkingLotUseCase(ctrl)
	r, wsHub := setupParkingLotHandler(mockUseCase)
	defer wsHub.Stop()

	mockUseCase.EXPECT().GetParkingLotWithOwnership(uint(1), "auth0|6721363081b8547d3f95a976").Return(&usecase.ParkingLotResponse{}, nil).Times(0)
	mockUseCase.EXPECT().UpdateParkingLot(uint(1), gomock.Any(), "auth0|6721363081b8547d3f95a976").Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(usecase.UpdateParkingLotRequest{
		Name:      "Updated Lot",
		Address:   "New Address",
		Latitude:  1.23,
		Longitude: 4.56,
	})
	req, _ := http.NewRequest("PUT", "/parkinglots/1", bytes.NewBuffer(body))
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "parking lot updated", response["status"])
}

func TestDeleteParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIParkingLotUseCase(ctrl)
	r, wsHub := setupParkingLotHandler(mockUseCase)
	defer wsHub.Stop()

	mockUseCase.EXPECT().GetParkingLotWithOwnership(uint(1), "auth0|6721363081b8547d3f95a976").Return(&usecase.ParkingLotResponse{}, nil).Times(0)
	mockUseCase.EXPECT().DeleteParkingLot(uint(1), "auth0|6721363081b8547d3f95a976").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/parkinglots/1", nil)
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response gin.H
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "parking lot deleted", response["status"])
}

func TestListParkingLots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mockgen.NewMockIParkingLotUseCase(ctrl)
	r, wsHub := setupParkingLotHandler(mockUseCase)
	defer wsHub.Stop()

	mockParkingLots := []usecase.ParkingLotResponse{
		{ID: 1, Name: "Lot 1", Address: "123 Test St", Latitude: 12.34, Longitude: 56.78, AvailableSpaces: 10},
		{ID: 2, Name: "Lot 2", Address: "456 Test Ave", Latitude: 98.76, Longitude: 54.32, AvailableSpaces: 20},
	}
	mockUseCase.EXPECT().ListParkingLots().Return(mockParkingLots, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/parkinglots/", nil)
	tokenString, _ := generateTestJWT()
	req.Header.Set("Authorization", "Bearer "+tokenString)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseParkingLots []usecase.ParkingLotResponse
	json.Unmarshal(w.Body.Bytes(), &responseParkingLots)
	assert.Equal(t, mockParkingLots, responseParkingLots)
}
