package domain

import (
	"time"

	"gorm.io/gorm"
)

type ParkingLot struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"not null" json:"name"`
	Latitude        float64        `gorm:"not null" json:"latitude"`
	Longitude       float64        `gorm:"not null" json:"longitude"`
	TotalSpaces     int            `gorm:"not null" json:"total_spaces"`
	AvailableSpaces int            `gorm:"not null" json:"available_spaces"`
	Status          string         `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
