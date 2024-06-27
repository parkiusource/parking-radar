package domain

import (
	"time"

	"gorm.io/gorm"
)

type ParkingLot struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"not null" json:"name"`
	Location        string         `gorm:"type:json;not null" json:"location"`
	TotalSpaces     int            `gorm:"not null" json:"total_spaces"`
	AvailableSpaces int            `gorm:"not null" json:"available_spaces"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
