package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	EXTERNAL_EVENT_PROVIDER_MICROSOFT ExternalEventProvider = "microsoft"
)

type ExternalEventProvider string

type Absence struct {
	gorm.Model
	UserID                *uint `gorm:"not null"`
	User                  *User `json:"-"`
	AbsenceFrom           time.Time
	AbsenceTill           time.Time
	AbsenceReasonID       *uint `gorm:"not null"`
	AbsenceReason         AbsenceReason
	SignedUserID          *uint
	SignedUser            *User
	ExternalEventProvider ExternalEventProvider
	ExternalEventID       string
	Identifier            uuid.UUID
	ExternalEvents        []AbsenceExternalEvent
}

type AbsenceExternalEvent struct {
	gorm.Model
	AbsenceID             uint
	Absence               Absence
	ExternalEventID       string
	ExternalEventProvider ExternalEventProvider
	Update                bool
}

type AbsenceReason struct {
	gorm.Model
	Description string
}

type AbsenceReasonCreateRequest struct {
	Description string `binding:"required"`
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

type AbsenceReturn struct {
	ID              uint
	User            UserResponse
	AbsenceFrom     time.Time
	AbsenceTill     time.Time
	AbsenceReasonID uint `json:",omitempty"`
	NettoDays       int
	CreatedAt       time.Time
	Reason          string `json:",omitempty"`
	Deletable       bool
}

func (a *Absence) GetAbsenceWorkDays() int {
	days := 0

	currentDay := a.AbsenceFrom

	for !currentDay.After(a.AbsenceTill) {
		// TODO: add holidays
		if currentDay.Weekday() != time.Saturday && currentDay.Weekday() != time.Sunday {
			days++
		}

		currentDay = currentDay.Add(24 * time.Hour)
	}

	return days
}

func (a *Absence) IsDeletableByUser() bool {
	return a.AbsenceFrom.After(time.Now()) || time.Now().Sub(a.CreatedAt).Hours() <= 24
}

func AbsenceReturns(absences []Absence, user *User, withReason bool, showRealReason bool) []AbsenceReturn {
	result := []AbsenceReturn{}
	for _, absence := range absences {
		returnObj := AbsenceReturn{
			ID:          absence.ID,
			AbsenceFrom: absence.AbsenceFrom,
			AbsenceTill: absence.AbsenceTill,
			NettoDays:   absence.GetAbsenceWorkDays(),
			CreatedAt:   absence.CreatedAt,
			Deletable:   absence.IsDeletableByUser(),
		}

		if user != nil {
			returnObj.User = user.GetUserResponse()
		} else {
			returnObj.User = absence.User.GetUserResponse()
		}

		if withReason {
			if showRealReason {
				returnObj.Reason = absence.AbsenceReason.Description
			} else {
				returnObj.Reason = "Abwesend"
			}
		} else {

			returnObj.AbsenceReasonID = *absence.AbsenceReasonID
		}

		result = append(result, returnObj)
	}

	return result
}
