package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"
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

	var result ExternalWorkCompensationHourSlots
	err := json.Unmarshal(bytes, &result)
	*e = result
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

	var result ExternalWorkCompensationAdditionalOptions
	err := json.Unmarshal(bytes, &result)
	*e = result
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

func (e ExternalWorkCompensationAdditionalOptions) Keys() []string {
	allOptions := []string{}
	for key, _ := range e {
		allOptions = append(allOptions, key)
	}

	sort.Strings(allOptions)
	return allOptions
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

const (
	EXTERNAL_WORK_STATUS_PLANNED      ExternalWorkStatus = "planned"
	EXTERNAL_WORK_STATUS_MISSING_INFO ExternalWorkStatus = "missing_info"
	EXTERNAL_WORK_STATUS_IN_REVIEW    ExternalWorkStatus = "in_review"
	EXTERNAL_WORK_STATUS_ACCEPTED     ExternalWorkStatus = "accepted"
	EXTERNAL_WORK_STATUS_DECLINED     ExternalWorkStatus = "declined"
	EXTERNAL_WORK_STATUS_INVOICED     ExternalWorkStatus = "invoiced"
)

type ExternalWorkStatus string

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
	ExternalWorkCompensationID uint
	ExternalWorkCompensation   ExternalWorkCompensation
	Status                     ExternalWorkStatus `gorm:"default:planned"`
	ReviewedBy                 *User
	ReviewedByID               *uint
	ReviewedDate               *time.Time
	InvoiceDate                *time.Time
	InvoiceIdentifier          *uuid.UUID `gorm:"index"`
}

func (e *ExternalWork) IsEditable() bool {
	return e.Status == EXTERNAL_WORK_STATUS_PLANNED || e.Status == EXTERNAL_WORK_STATUS_MISSING_INFO
}

func (e *ExternalWork) Calculate(holidays Holidays) ExternalWorkCalculated {
	externalWorkCalculated := ExternalWorkCalculated{
		ExternalWork:                        *e,
		TotalOvertimeHours:                  0,
		TotalExpensesWithoutSocialInsurance: 0,
		TotalExpensesWithSocialInsurance:    0,
		TotalOptions:                        make(map[string]float64),
		IsLocked:                            !e.IsEditable(),
	}

	externalWorkCalculated.WorkExpansesCalculated = []ExternalWorkExpenseCalculated{}

	for _, workExpense := range e.WorkExpanses {
		workExpense.ExternalWork = *e
		workExpenseCalculated := workExpense.Calculate(holidays)

		externalWorkCalculated.TotalExpensesWithSocialInsurance += workExpenseCalculated.ExpensesWithSocialInsurance
		externalWorkCalculated.TotalExpensesWithoutSocialInsurance += workExpenseCalculated.ExpensesWithoutSocialInsurance
		externalWorkCalculated.TotalOvertimeHours += workExpenseCalculated.TotalOvertimeHours

		externalWorkCalculated.WorkExpansesCalculated = append(externalWorkCalculated.WorkExpansesCalculated, workExpenseCalculated)

		for key, value := range workExpenseCalculated.AdditionalOptionsUsed {
			if _, exists := externalWorkCalculated.TotalOptions[key]; !exists {
				externalWorkCalculated.TotalOptions[key] = 0.0
			}

			externalWorkCalculated.TotalOptions[key] += value
		}
	}

	return externalWorkCalculated
}

type ExternalWorkCalculated struct {
	ExternalWork
	TotalOvertimeHours                  float64
	TotalExpensesWithoutSocialInsurance float64
	TotalExpensesWithSocialInsurance    float64
	TotalOptions                        map[string]float64
	IsLocked                            bool
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
	AdditionalOptions      StringArray `gorm:"type:jsonb" sql:"json"`
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

	additionalOptionsUsed := ExternalWorkCompensationAdditionalOptions{}

	for _, key := range e.ExternalWork.ExternalWorkCompensation.AdditionalOptions.Keys() {
		value := e.ExternalWork.ExternalWorkCompensation.AdditionalOptions[key]
		if slices.Contains(e.AdditionalOptions, key) {
			additionalOptionsUsed[key] = value
		} else {
			additionalOptionsUsed[key] = 0.0
		}
	}
	log.Printf("%#v", additionalOptionsUsed)

	totalOperationHours := e.TravelDurationHours + e.RestDurationHours + totalWorkingHours

	return ExternalWorkExpenseCalculated{
		ExternalWorkExpense:            *e,
		TotalOperationHours:            totalOperationHours,
		TotalWorkingHours:              totalWorkingHours,
		TotalOvertimeHours:             totalOperationHours - neededWorkingHours,
		TotalAwayHours:                 totalHoursAway,
		ExpensesWithoutSocialInsurance: expensesWithoutSocialInsurance,
		ExpensesWithSocialInsurance:    expensesWithSocialInsurance,
		AdditionalOptionsUsed:          additionalOptionsUsed,
		TravelPrivateKmCosts:           e.TravelWithPrivateCarKm * e.ExternalWork.ExternalWorkCompensation.PrivateCarKmCompensation,
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
	AdditionalOptionsUsed          ExternalWorkCompensationAdditionalOptions
	TravelPrivateKmCosts           float64
}
type DayDate struct {
	time.Time
}

func (t DayDate) MarshalJSON() ([]byte, error) {
	date := t.Time.Format("2006-01-02")
	date = fmt.Sprintf(`"%s"`, date)
	return []byte(date), nil
}

func (t *DayDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

type ExternalWorkCreateRequest struct {
	From                       DayDate `binding:"required,ltefield=Till" time_format:"2006-01-02"`
	Till                       DayDate `binding:"required" time_format:"2006-01-02"`
	Description                string  `binding:"required"`
	ExternalWorkCompensationID uint    `binding:"required"`
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
	AdditionalOptions      []string
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
	AdditionalOptions      []string
}

type ExternalWorkExternalEvent struct {
	gorm.Model
	ExternalWorkID        uint
	ExternalWork          ExternalWork
	ExternalEventID       string
	ExternalEventProvider ExternalEventProvider
	Update                bool
}
