package usecase

import (
	"errors"
	"testing"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/test/shared/mockgen"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Test for creating a parking lot.
func TestCreateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	adminRep := mockgen.NewMockIAdminRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo, adminRep)

	req := CreateParkingLotRequest{
		Name:      "Test Lot",
		Address:   "123 Test St",
		Latitude:  40.7128,
		Longitude: -74.0060,
		AdminUUID: "admin123",
	}

	adminRep.EXPECT().FindByAuth0UUID(req.AdminUUID).Return(&domain.Admin{ID: 123}, nil)

	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	_, err := useCase.CreateParkingLot(req)
	assert.NoError(t, err)
}

// Test for getting a parking lot with ownership.
func TestGetParkingLotWithOwnership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	adminRep := mockgen.NewMockIAdminRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo, adminRep)

	adminID := "admin123"
	parkingLotID := uint(1)

	adminRep.EXPECT().FindByAuth0UUID(adminID).Return(&domain.Admin{ID: 123}, nil)
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

// Test for getting a parking lot with ownership.
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

// Test for listing parking lots.
func TestListParkingLots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	adminRep := mockgen.NewMockIAdminRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo, adminRep)

	mockRepo.EXPECT().List().Return([]domain.ParkingLot{
		{ID: 1, Name: "Lot 1", Address: "Address 1", Latitude: 1.0, Longitude: 1.0},
	}, nil)

	sensorRepo.EXPECT().ListGroupedByParkingLot().Return(map[uint]uint{
		1: 2,
	}, nil)

	response, err := useCase.ListParkingLots()
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, uint(2), response[0].AvailableSpaces)
}

// Test for updating a parking lot.
func TestUpdateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	adminRep := mockgen.NewMockIAdminRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo, adminRep)

	adminID := "admin123"
	parkingLotID := uint(1)

	adminRep.EXPECT().FindByAuth0UUID(adminID).Return(&domain.Admin{ID: 123}, nil)

	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, uint(123)).Return(&domain.ParkingLot{
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

// Test for deleting a parking lot.
func TestDeleteParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	adminRep := mockgen.NewMockIAdminRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo, adminRep)

	adminID := "admin123"
	parkingLotID := uint(1)

	adminRep.EXPECT().FindByAuth0UUID(adminID).Return(&domain.Admin{ID: 123}, nil)

	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, uint(123)).Return(&domain.ParkingLot{
		ID: 1, AdminID: 123,
	}, nil)

	mockRepo.EXPECT().Delete(parkingLotID).Return(nil)

	err := useCase.DeleteParkingLot(parkingLotID, adminID)
	assert.NoError(t, err)
}

// Test for updating a parking lot with no ownership.
func TestUpdateParkingLot_NoOwnership(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	sensorRepo := mockgen.NewMockISensorRepository(ctrl)
	adminRep := mockgen.NewMockIAdminRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo, sensorRepo, adminRep)

	adminID := "wrong-admin"
	parkingLotID := uint(1)

	adminRep.EXPECT().FindByAuth0UUID(adminID).Return(&domain.Admin{ID: 123}, nil)

	mockRepo.EXPECT().GetByIDWithAdmin(parkingLotID, uint(123)).Return(nil, errors.New("not authorized"))

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
