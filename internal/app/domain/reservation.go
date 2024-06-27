package domain

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"not null" json:"user_id"`
	User             User           `gorm:"foreignKey:UserID" json:"user"`
	ParkingLotID     uint           `gorm:"not null" json:"parking_lot_id"`
	ParkingLot       ParkingLot     `gorm:"foreignKey:ParkingLotID" json:"parking_lot"`
	ReservationStart time.Time      `gorm:"not null" json:"reservation_start"`
	ReservationEnd   time.Time      `gorm:"not null" json:"reservation_end"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
