package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	EXTERNAL_EVENT_PROVIDER_MICROSOFT ExternalEventProvider = "microsoft"
	SIGNED_STATUS_ACCEPTED            AbsenceSignedStatus   = "accepted"
	SIGNED_STATUS_DECLINED            AbsenceSignedStatus   = "declined"
)

type AbsenceSignedStatus string
type ExternalEventProvider string

type Absence struct {
	gorm.Model
	UserID                *uint `gorm:"not null"`
	User                  *User
	AbsenceFrom           time.Time
	AbsenceTill           time.Time
	AbsenceReasonID       *uint `gorm:"not null"`
	AbsenceReason         AbsenceReason
	SignedUserID          *uint
	SignedUser            *User
	SignedMessage         *string
	SignedStatus          *AbsenceSignedStatus
	SignedTimestamp       *time.Time
	ExternalEventProvider ExternalEventProvider
	ExternalEventID       string
	Identifier            uuid.UUID
	ExternalEvents        []AbsenceExternalEvent
	NettoDays             *float64
}

type AbsenceExternalEvent struct {
	gorm.Model
	AbsenceID             uint
	Absence               Absence
	ExternalEventID       string
	ExternalEventProvider ExternalEventProvider
	Update                bool
}

const (
	ABESENCE_REASON_OVERTIME_IMPACT_NONE     AbsenceReasonOvertimeImpact = "none"
	ABESENCE_REASON_OVERTIME_IMPACT_DURATION AbsenceReasonOvertimeImpact = "duration"
	ABESENCE_REASON_OVERTIME_IMPACT_HOURS    AbsenceReasonOvertimeImpact = "hours"
	ABESENCE_REASON_OVERTIME_IMPACT_DAYS     AbsenceReasonOvertimeImpact = "days"
)

type AbsenceReasonOvertimeImpact string

type AbsenceReason struct {
	gorm.Model
	Description    string
	OvertimeImpact AbsenceReasonOvertimeImpact `gorm:"default:none"`
	ImpactHours    float64
	ImpactDays     float64
	NeedsApproval  *bool
}

type AbsenceReasonCreateRequest struct {
	Description    string `binding:"required"`
	NeedsApproval  *bool
	OvertimeImpact AbsenceReasonOvertimeImpact
	ImpactHours    float64
	ImpactDays     float64
}

type AbsenceSignRequest struct {
	Status  AbsenceSignedStatus `binding:"required"`
	Message *string
}

type AbsenceCreateRequest struct {
	AbsenceFrom     string `binding:"required" time_format:"2006-01-02"`
	AbsenceTill     string `binding:"required" time_format:"2006-01-02"`
	AbsenceReasonID uint   `binding:"required"`
}

func (acr *AbsenceCreateRequest) AbsenceFromParsed() (time.Time, error) {
	return time.Parse("2006-01-02", acr.AbsenceFrom)
}

func (acr *AbsenceCreateRequest) AbsenceTillParsed() (time.Time, error) {
	return time.Parse("2006-01-02", acr.AbsenceTill)
}

type AbsenceUserSummaryYearReason struct {
	Upcoming float64
	Past     float64
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
	AbsenceReason   AbsenceReason
	SignedUserID    *uint
	SignedUser      *User
	SignedMessage   *string
	SignedStatus    *AbsenceSignedStatus
	SignedTimestamp *time.Time
	NettoDays       float64
	CreatedAt       time.Time
	Reason          string `json:",omitempty"`
	Deletable       bool
}

func (a *Absence) CalculateNettoDays(holidays Holidays) {
	days := 0

	currentDay := a.AbsenceFrom

	for !currentDay.After(a.AbsenceTill) {
		if currentDay.Weekday() != time.Saturday && currentDay.Weekday() != time.Sunday && !holidays.Contains(currentDay) {
			days++
		}

		currentDay = currentDay.Add(24 * time.Hour)
	}

	total := float64(days)

	a.NettoDays = &total
}

func (a *Absence) IsDeletableByUser() bool {
	return a.AbsenceFrom.After(time.Now()) || time.Now().Sub(a.CreatedAt).Hours() <= 24
}

func (a *Absence) IsDateInAbsence(d time.Time) bool {
	return d.Equal(a.AbsenceFrom.UTC()) || d.Equal(a.AbsenceTill.UTC()) || (d.Before(a.AbsenceTill.UTC()) && d.After(a.AbsenceFrom.UTC()))
}

func (a *Absence) Sign(signingUser *User, status AbsenceSignedStatus, message *string) {
	now := time.Now()

	a.SignedUser = signingUser
	a.SignedTimestamp = &now
	a.SignedStatus = &status
	a.SignedMessage = message
}

func AbsenceReturns(absences []Absence, user *User, withReason bool, showRealReason bool, withSigningInfo bool) []AbsenceReturn {
	result := []AbsenceReturn{}
	for _, absence := range absences {
		returnObj := AbsenceReturn{
			ID:          absence.ID,
			AbsenceFrom: absence.AbsenceFrom,
			AbsenceTill: absence.AbsenceTill,
			NettoDays:   *absence.NettoDays,
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

		if withSigningInfo {
			returnObj.SignedUserID = absence.SignedUserID
			returnObj.SignedMessage = absence.SignedMessage
			returnObj.SignedStatus = absence.SignedStatus
			returnObj.SignedTimestamp = absence.SignedTimestamp
		}

		result = append(result, returnObj)
	}

	return result
}
