package domain

import "time"

type Esp32Device struct {
	ID                uint64    `gorm:"primaryKey" json:"id"`
	DeviceIdentifier  string    `json:"device_identifier"`
	LastCommunication time.Time `json:"last_communication"`
}
