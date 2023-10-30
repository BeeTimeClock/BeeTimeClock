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

type TimestampMonthQuota struct {
	gorm.Model
	UserID uint `gorm:"index:idx_month_quota,unique;index"`
	User   User
	Year   int `gorm:"index:idx_month_quota,unique"`
	Month  int `gorm:"index:idx_month_quota,unique"`
	Hours  float64
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
	ChangeReason    string `binding:"required"`
}

type TimestampActionCheckInRequest struct {
	IsHomeoffice bool
}

type TimestampCorrectionCreateRequest struct {
	ChangeReason       string    `binding:"required"`
	NewComingTimestamp time.Time `binding:"required"`
	NewGoingTimestamp  time.Time `binding:"required"`
	IsHomeoffice       bool
}

type TimestampGroup struct {
	Date            time.Time
	IsHomeoffice    bool
	Timestamps      []Timestamp
	WorkingHours    float64
	SubtractedHours float64
	OvertimeHours   float64
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

type WorkTimeModel struct {
	DefaultHoursPerWeekday   float64
	HoursPerWeekdayException map[time.Weekday]float64
}

func DefaultWorkTimeModel() WorkTimeModel {
	return WorkTimeModel{
		DefaultHoursPerWeekday: 8.0,
		HoursPerWeekdayException: map[time.Weekday]float64{
			time.Friday: 6.0,
		},
	}
}
