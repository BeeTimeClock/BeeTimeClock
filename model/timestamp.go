package model

import (
	"time"

	"gorm.io/gorm"
)

type Timestamp struct {
	gorm.Model
	UserID          uint  `gorm:"not null"`
	User            *User `json:"-"`
	ComingTimestamp time.Time
	GoingTimestamp  time.Time
	IsHomeoffice    bool
	Corrections     []TimestampCorrection
}

type TimestampCorrection struct {
	gorm.Model
	TimestampID        uint `gorm:"not null"`
	Timestamp          Timestamp
	ChangeReason       string
	OldComingTimestamp time.Time
	OldGoingTimestamp  time.Time
}

type TimestampCreateRequest struct {
	ComingTimestamp time.Time `binding:"required"`
	GoingTimestamp  time.Time
	IsHomeoffice    bool
}

type TimestampActionCheckInRequest struct {
	IsHomeoffice bool
}

type TimestampCorrectionCreateRequest struct {
	ChangeReason       string    `binding:"required"`
	NewComingTimestamp time.Time `binding:"required"`
	NewGoingTimestamp  time.Time `binding:"required"`
}

func (t *Timestamp) IsComplete() bool {
	return !t.GoingTimestamp.IsZero()
}

func (t *Timestamp) CalculateWorkingHours() (float64, float64) {
	goingTimestamp := t.GoingTimestamp

	if goingTimestamp.IsZero() {
		goingTimestamp = time.Now()
	}

	completeTime := goingTimestamp.Sub(t.ComingTimestamp).Hours()
	calculatedTime := completeTime

	if completeTime > 6 {
		calculatedTime = calculatedTime - 0.5
	}

	if completeTime > 9 {
		calculatedTime = calculatedTime - 0.25
	}

	return calculatedTime, completeTime - calculatedTime
}
