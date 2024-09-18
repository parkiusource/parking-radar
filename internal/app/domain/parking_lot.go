package domain

import (
	"time"

	"gorm.io/gorm"
)

type ParkingLot struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Address      string         `gorm:"type:varchar(40);not null" json:"address"`
	Latitude     float64        `gorm:"not null" json:"latitude"`
	Longitude    float64        `gorm:"not null" json:"longitude"`
	ContactName  string         `gorm:"type:varchar(40);not null" json:"contact_name"`
	ContactPhone string         `gorm:"type:varchar(40);not null" json:"contact_phone"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
