package model

import "gorm.io/gorm"

type Settings struct {
	gorm.Model

	CheckinDetectionByIPAddress        *bool                       `gorm:"default:false"`
	OfficeIPAddresses                  []SettingsOfficeIPAddresses `gorm:"constraint:OnDelete:CASCADE"`
	TimestampChangeReasonMinimumLength int64                       `gorm:"default:20"`
}

type SettingsOfficeIPAddresses struct {
	gorm.Model
	Settings    Settings
	SettingsID  uint
	IPAddress   string `gorm:"uniqueIndex"`
	Description string
}
