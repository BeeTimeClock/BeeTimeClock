package model

import (
	"time"

	"gorm.io/gorm"
)

type Absence struct {
	gorm.Model
	UserID          *uint `gorm:"not null"`
	User            *User `json:"-"`
	AbsenceFrom     time.Time
	AbsenceTill     time.Time
	AbsenceReasonID *uint `gorm:"not null"`
	AbsenceReason   AbsenceReason
	SignedUserID    *uint
	SignedUser      *User
}

type AbsenceReason struct {
	gorm.Model
	Description string
}

type AbsenceCreateRequest struct {
	AbsenceFrom     time.Time `binding:"required"`
	AbsenceTill     time.Time `binding:"required"`
	AbsenceReasonID uint      `binding:"required"`
}

type AbsenceUserSummaryYearReason struct {
	Upcoming int
	Past     int
}

type AbsenceUserSummaryYear struct {
	ByAbsenceReason map[uint]AbsenceUserSummaryYearReason
}

type AbsenceUserSummary struct {
	ByYear             map[int]AbsenceUserSummaryYear
	HolidayDaysPerYear uint
}
