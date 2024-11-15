package domain

import "time"

type Admin struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	Auth0UUID    string       `gorm:"uniqueIndex" json:"auth0_uuid"`
	NIT          string       `gorm:"not null" json:"nit"`
	PhotoURL     string       `json:"photo_url"`
	ContactPhone string       `json:"contact_phone"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	ParkingLots  []ParkingLot `gorm:"foreignKey:AdminID" json:"parking_lots"`
}

type AdminProfileData struct {
	NIT          string `json:"nit"`
	PhotoURL     string `json:"photo_url"`
	ContactPhone string `json:"contact_phone"`
}
