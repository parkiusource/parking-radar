package domain

import (
	"time"

	"gorm.io/gorm"
)

type Sensor struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Esp32DeviceID uint           `gorm:"not null" json:"esp32_device_id"`
	Esp32Device   Esp32Device    `gorm:"foreignKey:Esp32DeviceID" json:"esp32_device"`
	ParkingLotID  uint           `gorm:"not null" json:"parking_lot_id"`
	ParkingLot    ParkingLot     `gorm:"foreignKey:ParkingLotID" json:"parking_lot"`
	Status        string         `gorm:"not null" json:"status"`
	SensorNumber  int            `json:"sensor_number"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
