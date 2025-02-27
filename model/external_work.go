package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExternalWork struct {
	gorm.Model
	UserID       uint
	User         User
	Description  string
	From         time.Time
	Till         time.Time
	Identifier   uuid.UUID
	WorkExpanses []ExternalWorkExpense
	ClosedTime   *time.Time
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
	OnSiteFrom          *time.Time
	OnSiteTill          *time.Time
	Place               string
}

type ExternalWorkCreateRequest struct {
	From        time.Time `binding:"required"`
	Till        time.Time `binding:"required"`
	Description string    `binding:"required"`
}

type ExternalWorkExternalEvent struct {
	gorm.Model
	ExternalWorkID        uint
	ExternalWork          ExternalWork
	ExternalEventID       string
	ExternalEventProvider ExternalEventProvider
	Update                bool
}
