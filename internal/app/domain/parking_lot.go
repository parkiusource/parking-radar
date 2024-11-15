package domain

import (
	"time"

	"gorm.io/gorm"
)

type ParkingLot struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Address      string         `gorm:"type:varchar(40);not null" json:"address"`
	Latitude     float64        `gorm:"not null;uniqueIndex:idx_lat_long" json:"latitude"`
	Longitude    float64        `gorm:"not null;uniqueIndex:idx_lat_long" json:"longitude"`
	ContactName  string         `gorm:"type:varchar(40);not null" json:"contact_name"`
	ContactPhone string         `gorm:"type:varchar(40);not null" json:"contact_phone"`
	AdminID      uint           `gorm:"not null" json:"admin_id"`
	Admin        Admin          `gorm:"foreignKey:AdminID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
