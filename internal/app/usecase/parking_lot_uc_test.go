package usecase

import (
	"testing"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/test/shared/mockgen"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Helper to set up common dependencies for tests.
func setupTest(t *testing.T) (*gomock.Controller, *mockgen.MockIParkingLotRepository, *mockgen.MockISensorRepository, *mockgen.MockIAdminRepository, IParkingLotUseCase) {
	ctrl := gomock.NewController(t)
	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	adminRepo := mockgen.NewMockIAdminRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo, adminRepo)
	return ctrl, mockRepo, sensorRepo, adminRepo, useCase
}

func TestCreateParkingLot(t *testing.T) {
	ctrl, mockRepo, _, adminRepo, useCase := setupTest(t)
	defer ctrl.Finish()

	req := CreateParkingLotRequest{
		Name:      "Test Lot",
		Address:   "123 Test St",
		Latitude:  40.7128,
		Longitude: -74.0060,
		AdminUUID: "admin123",
	}

	adminRepo.EXPECT().FindByAuth0UUID(req.AdminUUID).Return(&domain.Admin{ID: 123}, nil)
	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	_, err := useCase.CreateParkingLot(req)
	assert.NoError(t, err)
}

func TestGetParkingLotWithOwnership(t *testing.T) {
	ctrl, mockRepo, sensorRepo, adminRepo, useCase := setupTest(t)
	defer ctrl.Finish()

	adminID := "admin123"
	parkingLotID := uint(1)

	adminRepo.EXPECT().FindByAuth0UUID(adminID).Return(&domain.Admin{ID: 123}, nil)
	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, uint(123)).Return(&domain.ParkingLot{
		ID: 1, Name: "Test Lot", Address: "123 Test St", Latitude: 40.7128, Longitude: -74.0060,
	}, nil)
	sensorRepo.EXPECT().ListByParkingLot(parkingLotID).Return([]domain.Sensor{
		{Status: "free"},
		{Status: "busy"},
	}, nil)

	response, err := useCase.GetParkingLotWithOwnership(parkingLotID, adminID)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, uint(1), response.AvailableSpaces)
}

func TestGetParkingLot(t *testing.T) {
	ctrl, mockRepo, sensorRepo, _, useCase := setupTest(t)
	defer ctrl.Finish()

	parkingLotID := uint(1)

	mockRepo.EXPECT().GetByID(parkingLotID).Return(&domain.ParkingLot{
		ID: 1, Name: "Test Lot", Address: "123 Test St", Latitude: 40.7128, Longitude: -74.0060,
	}, nil)
	sensorRepo.EXPECT().ListByParkingLot(parkingLotID).Return([]domain.Sensor{
		{Status: "free"},
		{Status: "free"},
	}, nil)

	response, err := useCase.GetParkingLot(parkingLotID)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, uint(2), response.AvailableSpaces)
}
