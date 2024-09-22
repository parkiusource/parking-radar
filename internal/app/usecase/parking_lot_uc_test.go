package usecase

import (
	"github.com/CamiloLeonP/parking-radar/internal/test/shared/mockgen"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
)

func TestCreateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	useCase := NewParkingLotUseCase(mockRepo)

	req := CreateParkingLotRequest{
		Name:      "test Lot",
		Address:   "123 test St",
		Latitude:  40.7128,
		Longitude: -74.0060,
	}

	mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

	err := useCase.CreateParkingLot(req)
	assert.NoError(t, err)
}

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

func TestListParkingLots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	mockSensorRepo := mockgen.NewMockISensorRepository(ctrl)
	useCase := &ParkingLotUseCase{
		ParkingLotRepository: mockRepo,
		SensorRepository:     mockSensorRepo,
	}

	mockRepo.EXPECT().List().Return([]domain.ParkingLot{
		{ID: 1, Name: "Lot 1", Address: "Address 1", Latitude: 1.0, Longitude: 1.0},
	}, nil)

	mockSensorRepo.EXPECT().ListByParkingLot(uint(1)).Return([]domain.Sensor{
		{Status: "free"},
		{Status: "free"},
	}, nil)

	response, err := useCase.ListParkingLots()
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, uint(2), response[0].AvailableSpaces)
}

func TestUpdateParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	useCase := &ParkingLotUseCase{
		ParkingLotRepository: mockRepo,
	}

	mockRepo.EXPECT().GetByID(uint(1)).Return(&domain.ParkingLot{
		ID: 1, Name: "Old Name", Address: "Old Address", Latitude: 1.0, Longitude: 1.0,
	}, nil)

	mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

	req := UpdateParkingLotRequest{
		Name:      "New Name",
		Address:   "New Address",
		Latitude:  2.0,
		Longitude: 2.0,
	}

	err := useCase.UpdateParkingLot(1, req)
	assert.NoError(t, err)
}

func TestDeleteParkingLot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockgen.NewMockIParkingLotRepository(ctrl)
	useCase := &ParkingLotUseCase{
		ParkingLotRepository: mockRepo,
	}

	mockRepo.EXPECT().Delete(uint(1)).Return(nil)

	err := useCase.DeleteParkingLot(1)
	assert.NoError(t, err)
}
