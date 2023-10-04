package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	FUEL_STATE_PREPARED FuelState = "prepared"
	FUEL_STATE_OPEN     FuelState = "open"
	FUEL_STATE_EXPORTED FuelState = "exported"
)

type FuelState string

type Fuel struct {
	gorm.Model
	UserID           *uint `gorm:"not null"`
	User             *User `json:"-"`
	ReceiptDate      time.Time
	ReceiptValue     float64
	ReceiptFuelValue float64
	ReceiptRawText   string
	ExportedAt       *time.Time
	UploadFileName   string
	State            FuelState
}

type FuelUpdateRequest struct {
	ReceiptDate      time.Time `binding:"required"`
	ReceiptValue     float64   `binding:"required"`
	ReceiptFuelValue float64   `binding:"required"`
}

type FuelQueryAll struct {
	State FuelState
}

func (fq *FuelQueryAll) HasState() bool {
	return fq.State != ""
}
