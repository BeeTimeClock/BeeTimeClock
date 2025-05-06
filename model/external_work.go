package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExternalWorkCompensationHourSlot struct {
	Hours        float64
	Compensation float64
}

type ExternalWorkCompensationHourSlots []ExternalWorkCompensationHourSlot

func (e *ExternalWorkCompensationHourSlots) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := &ExternalWorkCompensationHourSlots{}
	err := json.Unmarshal(bytes, &result)
	*e = ExternalWorkCompensationHourSlots(*result)
	return err
}

func (e ExternalWorkCompensationHourSlots) Value() (driver.Value, error) {
	if len(e) == 0 {
		return nil, nil
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

type ExternalWorkCompensationAdditionalOptions map[string]float64

func (e *ExternalWorkCompensationAdditionalOptions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := &ExternalWorkCompensationAdditionalOptions{}
	err := json.Unmarshal(bytes, &result)
	*e = ExternalWorkCompensationAdditionalOptions(*result)
	return err
}

func (e ExternalWorkCompensationAdditionalOptions) Value() (driver.Value, error) {
	if len(e) == 0 {
		return nil, nil
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

type ExternalWorkCompensation struct {
	gorm.Model
	IsoCountryCodeA2            string                                    `gorm:"index:idx_compensation,unique;index"`
	WithSocialInsuranceSlots    ExternalWorkCompensationHourSlots         `gorm:"type:jsonb" sql:"json"`
	WithoutSocialInsuranceSlots ExternalWorkCompensationHourSlots         `gorm:"type:jsonb" sql:"json"`
	AdditionalOptions           ExternalWorkCompensationAdditionalOptions `gorm:"type:jsonb" sql:"json"`
	ValidFrom                   time.Time                                 `gorm:"index:idx_compensation,unique;index"`
	ValidTill                   time.Time                                 `gorm:"index:idx_compensation,unique;index"`
	PrivateCarKmCompensation    float64
}

type ExternalWork struct {
	gorm.Model
	UserID                     uint
	User                       User
	Description                string
	From                       time.Time
	Till                       time.Time
	Identifier                 uuid.UUID
	WorkExpanses               []ExternalWorkExpense           `json:"-"`
	WorkExpansesCalculated     []ExternalWorkExpenseCalculated `gorm:"-" json:"WorkExpanses"`
	ClosedTime                 *time.Time
	ExternalWorkCompensationID uint
	ExternalWorkCompensation   ExternalWorkCompensation
}

func (e *ExternalWork) Calculate(holidays Holidays) ExternalWorkCalculated {
	externalWorkCalculated := ExternalWorkCalculated{
		ExternalWork:                        *e,
		TotalOvertimeHours:                  0,
		TotalExpensesWithoutSocialInsurance: 0,
		TotalExpensesWithSocialInsurance:    0,
	}

	externalWorkCalculated.WorkExpansesCalculated = []ExternalWorkExpenseCalculated{}

	for _, workExpense := range e.WorkExpanses {
		workExpense.ExternalWork = *e
		workExpenseCalculated := workExpense.Calculate(holidays)

		externalWorkCalculated.TotalExpensesWithSocialInsurance += workExpenseCalculated.ExpensesWithSocialInsurance
		externalWorkCalculated.TotalExpensesWithoutSocialInsurance += workExpenseCalculated.ExpensesWithoutSocialInsurance
		externalWorkCalculated.TotalOvertimeHours += workExpenseCalculated.TotalOvertimeHours

		externalWorkCalculated.WorkExpansesCalculated = append(externalWorkCalculated.WorkExpansesCalculated, workExpenseCalculated)
	}

	return externalWorkCalculated
}

type ExternalWorkCalculated struct {
	ExternalWork
	TotalOvertimeHours                  float64
	TotalExpensesWithoutSocialInsurance float64
	TotalExpensesWithSocialInsurance    float64
}

type ExternalWorkExpense struct {
	gorm.Model
	ExternalWorkID         uint
	ExternalWork           ExternalWork
	Date                   time.Time
	DepartureTime          *time.Time
	ArrivalTime            *time.Time
	TravelDurationHours    float64
	PauseDurationHours     float64
	RestDurationHours      float64
	OnSiteFrom             *time.Time
	OnSiteTill             *time.Time
	Place                  string
	TravelWithPrivateCarKm float64
}

func (e *ExternalWorkExpense) Calculate(holidays Holidays) ExternalWorkExpenseCalculated {
	workTimeModel := DefaultWorkTimeModel()
	neededWorkingHours := workTimeModel.GetWorkingHoursForDay(e.Date, holidays)

	expensesWithSocialInsurance := 0.0
	expensesWithoutSocialInsurance := 0.0

	totalHoursAway := 24.0
	totalWorkingHours := 0.0
	startOfDay := e.Date.UTC().Round(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	if e.DepartureTime != nil && e.ArrivalTime == nil {
		_, offset := e.DepartureTime.Zone()
		totalHoursAway = endOfDay.Sub(e.DepartureTime.UTC()).Hours() - time.Duration(time.Second*time.Duration(offset)).Hours()
	}

	if e.DepartureTime == nil && e.ArrivalTime != nil {
		tz, offset := e.ArrivalTime.Zone()
		totalHoursAway = float64(e.ArrivalTime.Sub(startOfDay).Hours() - time.Duration(time.Second*time.Duration(offset)).Hours())
		log.Printf("only arrival %#v: %s: %f", e.ArrivalTime, tz, time.Duration(time.Second*time.Duration(offset)).Hours())
	}

	if e.DepartureTime != nil && e.ArrivalTime != nil {
		totalHoursAway = e.ArrivalTime.Sub(*e.DepartureTime).Hours()
	}

	if e.OnSiteFrom != nil && e.OnSiteTill != nil {
		totalWorkingHours = e.OnSiteTill.Sub(*e.OnSiteFrom).Hours()
	}

	for _, hourSlot := range e.ExternalWork.ExternalWorkCompensation.WithSocialInsuranceSlots {
		if totalHoursAway >= hourSlot.Hours {
			expensesWithSocialInsurance = hourSlot.Compensation
		}
	}

	for _, hourSlot := range e.ExternalWork.ExternalWorkCompensation.WithoutSocialInsuranceSlots {
		if totalHoursAway >= hourSlot.Hours {
			expensesWithoutSocialInsurance = hourSlot.Compensation
		}
	}

	totalOperationHours := e.TravelDurationHours + e.RestDurationHours + totalWorkingHours

	return ExternalWorkExpenseCalculated{
		ExternalWorkExpense:            *e,
		TotalOperationHours:            totalOperationHours,
		TotalWorkingHours:              totalWorkingHours,
		TotalOvertimeHours:             totalOperationHours - neededWorkingHours,
		TotalAwayHours:                 totalHoursAway,
		ExpensesWithoutSocialInsurance: expensesWithoutSocialInsurance,
		ExpensesWithSocialInsurance:    expensesWithSocialInsurance,
	}
}

type ExternalWorkExpenseCalculated struct {
	ExternalWorkExpense
	TotalOperationHours            float64
	TotalWorkingHours              float64
	TotalOvertimeHours             float64
	TotalAwayHours                 float64
	ExpensesWithoutSocialInsurance float64
	ExpensesWithSocialInsurance    float64
}

type ExternalWorkCreateRequest struct {
	From                       time.Time `binding:"required,ltefield=Till"`
	Till                       time.Time `binding:"required"`
	Description                string    `binding:"required"`
	ExternalWorkCompensationID uint      `binding:"required"`
}

type ExternalWorkExpenseCreateRequest struct {
	Date                   time.Time `binding:"required"`
	DepartureTime          *time.Time
	ArrivalTime            *time.Time
	TravelDurationHours    float64
	PauseDurationHours     float64
	OnSiteFrom             *time.Time
	OnSiteTill             *time.Time
	Place                  string
	TravelWithPrivateCarKm float64
}

type ExternalWorkExpenseUpdateRequest struct {
	DepartureTime          *time.Time
	ArrivalTime            *time.Time
	TravelDurationHours    float64
	PauseDurationHours     float64
	OnSiteFrom             *time.Time
	OnSiteTill             *time.Time
	Place                  string
	TravelWithPrivateCarKm float64
}

type ExternalWorkExternalEvent struct {
	gorm.Model
	ExternalWorkID        uint
	ExternalWork          ExternalWork
	ExternalEventID       string
	ExternalEventProvider ExternalEventProvider
	Update                bool
}
