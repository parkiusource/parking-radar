package usecase

import (
	"errors"
	"testing"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/test/shared/mockgen"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Simula el contexto con JWT y claims.
func mockContextWithClaims(adminID string) map[string]interface{} {
	return map[string]interface{}{
		"sub": adminID,
	}
}

// Test para crear un parking lot.
func TestCreateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo)

	req := CreateParkingLotRequest{
		Name:      "Test Lot",
		Address:   "123 Test St",
		Latitude:  40.7128,
		Longitude: -74.0060,
		AdminID:   "admin123",
	}

	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	_, err := useCase.CreateParkingLot(req)
	assert.NoError(t, err)
}

// Test para obtener un parking lot con ownership.
func TestGetParkingLotWithOwnership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	mockSensorRepo := mockgen.NewMockISensorRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, mockSensorRepo)

	adminID := "admin123"
	parkingLotID := uint(1)

	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, adminID).Return(&domain.ParkingLot{
		ID: 1, Name: "Test Lot", Address: "123 Test St", Latitude: 40.7128, Longitude: -74.0060,
	}, nil)

	mockSensorRepo.EXPECT().ListByParkingLot(parkingLotID).Return([]domain.Sensor{
		{Status: "free"},
		{Status: "busy"},
	}, nil)

	response, err := useCase.GetParkingLotWithOwnership(parkingLotID, adminID)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, uint(1), response.AvailableSpaces)
}

// Test para obtener un parking lot sin ownership.
func TestGetParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	mockSensorRepo := mockgen.NewMockISensorRepository(ctrl)
	useCase := &ParkingLotUseCase{
		ParkingLotRepository: mockRepo,
		SensorRepository:     mockSensorRepo,
	}

	mockRepo.EXPECT().GetByID(uint(1)).Return(&domain.ParkingLot{
		ID:        1,
		Name:      "test Lot",
		Address:   "123 test St",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}, nil)

	mockSensorRepo.EXPECT().ListByParkingLot(uint(1)).Return([]domain.Sensor{
		{Status: "free"},
		{Status: "occupied"},
	}, nil)

	response, err := useCase.GetParkingLot(1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, uint(1), response.AvailableSpaces)
}

// Test para listar parking lots.
func TestListParkingLots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	mockSensorRepo := mockgen.NewMockISensorRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, mockSensorRepo)

	mockRepo.EXPECT().List().Return([]domain.ParkingLot{
		{ID: 1, Name: "Lot 1", Address: "Address 1", Latitude: 1.0, Longitude: 1.0},
	}, nil)

	mockSensorRepo.EXPECT().ListGroupedByParkingLot().Return(map[uint]uint{
		1: 2,
	}, nil)

	response, err := useCase.ListParkingLots()
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, uint(2), response[0].AvailableSpaces)
}

// Test para actualizar un parking lot.
func TestUpdateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, nil)

	adminID := "admin123"
	parkingLotID := uint(1)

	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, adminID).Return(&domain.ParkingLot{
		ID: 1, Name: "Old Name", Address: "Old Address", Latitude: 1.0, Longitude: 1.0,
	}, nil)

	mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

	req := UpdateParkingLotRequest{
		Name:      "New Name",
		Address:   "New Address",
		Latitude:  2.0,
		Longitude: 2.0,
	}

	err := useCase.UpdateParkingLot(parkingLotID, req, adminID)
	assert.NoError(t, err)
}

// Test para eliminar un parking lot.
func TestDeleteParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, nil)

	adminID := "admin123"
	parkingLotID := uint(1)

	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, adminID).Return(&domain.ParkingLot{
		ID: 1, AdminID: adminID,
	}, nil)

	mockRepo.EXPECT().Delete(parkingLotID).Return(nil)

	err := useCase.DeleteParkingLot(parkingLotID, adminID)
	assert.NoError(t, err)
}

// Test para caso de error en actualización por falta de ownership.
func TestUpdateParkingLot_NoOwnership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, nil)

	adminID := "wrong-admin"
	parkingLotID := uint(1)

	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, adminID).Return(nil, errors.New("not authorized"))

	req := UpdateParkingLotRequest{
		Name:      "New Name",
		Address:   "New Address",
		Latitude:  2.0,
		Longitude: 2.0,
	}

	err := useCase.UpdateParkingLot(parkingLotID, req, adminID)
	assert.Error(t, err)
	assert.Equal(t, "not authorized", err.Error())
}
