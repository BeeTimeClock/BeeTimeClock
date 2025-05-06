package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExternalWork struct {
	gorm.Model
	UserID                 uint
	User                   User
	Description            string
	From                   time.Time
	Till                   time.Time
	Identifier             uuid.UUID
	WorkExpanses           []ExternalWorkExpense           `json:"-"`
	WorkExpansesCalculated []ExternalWorkExpenseCalculated `gorm:"-" json:"WorkExpanses"`
	ClosedTime             *time.Time
}

type ExternalWorkExpense struct {
	gorm.Model
	ExternalWorkID      uint
	ExternalWork        ExternalWork
	Date                time.Time
	DepartureTime       *time.Time
	ArrivalTime         *time.Time
	TravelDurationHours float64
	PauseDurationHours  float64
	RestDurationHours   float64
	OnSiteFrom          *time.Time
	OnSiteTill          *time.Time
	Place               string
}

func (e *ExternalWorkExpense) Calculate() ExternalWorkExpenseCalculated {
	workTimeModel := DefaultWorkTimeModel()
	neededWorkingHours := workTimeModel.GetWorkingHoursForDay(e.Date)

	totalHoursAway := 24.0
	totalWorkingHours := 0.0
	startOfDay := e.Date.UTC().Round(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	if e.DepartureTime != nil && e.ArrivalTime == nil {
		_, offset := e.DepartureTime.Zone()
		totalHoursAway = endOfDay.Sub(e.DepartureTime.UTC()).Hours() - time.Duration(time.Second*time.Duration(offset)).Hours()
	}

	if e.DepartureTime == nil && e.ArrivalTime != nil {
		_, offset := e.ArrivalTime.Zone()
		totalHoursAway = float64(e.ArrivalTime.Sub(startOfDay).Hours() - float64(offset))
	}

	if e.DepartureTime != nil && e.ArrivalTime != nil {
		totalHoursAway = e.ArrivalTime.Sub(*e.DepartureTime).Hours()
	}

	if e.OnSiteFrom != nil && e.OnSiteTill != nil {
		totalWorkingHours = e.OnSiteTill.Sub(*e.OnSiteFrom).Hours()
	}

	return ExternalWorkExpenseCalculated{
		ExternalWorkExpense: *e,
		TotalOperationHours: e.PauseDurationHours + e.TravelDurationHours + e.RestDurationHours + totalWorkingHours,
		TotalWorkingHours:   totalWorkingHours,
		TotalOvertimeHours:  totalWorkingHours - neededWorkingHours,
		TotalAwayHours:      totalHoursAway,
	}
}

type ExternalWorkExpenseCalculated struct {
	ExternalWorkExpense
	TotalOperationHours float64
	TotalWorkingHours   float64
	TotalOvertimeHours  float64
	TotalAwayHours      float64
}

type ExternalWorkCreateRequest struct {
	From        time.Time `binding:"required"`
	Till        time.Time `binding:"required"`
	Description string    `binding:"required"`
}

type ExternalWorkExpenseCreateRequest struct {
	Date                time.Time `binding:"required"`
	DepartureTime       *time.Time
	ArrivalTime         *time.Time
	TravelDurationHours float64
	PauseDurationHours  float64
	OnSiteFrom          *time.Time
	OnSiteTill          *time.Time
	Place               string
}

type ExternalWorkExternalEvent struct {
	gorm.Model
	ExternalWorkID        uint
	ExternalWork          ExternalWork
	ExternalEventID       string
	ExternalEventProvider ExternalEventProvider
	Update                bool
}
