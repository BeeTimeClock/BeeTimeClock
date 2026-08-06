package model

import (
	"gorm.io/gorm"
)

const (
	OVERTIME_SUBTRACTION_MODEL_HOURS      OvertimeSubtractionModel = "hours"
	OVERTIME_SUBTRACTION_MODEL_PERCENTAGE OvertimeSubtractionModel = "percentage"
)

type OvertimeSubtractionModel string

type WorkTimeModel struct {
	gorm.Model
	Name                      string `gorm:"unique"`
	WorkingHoursPerWeekday    float64
	OvertimeSubtractionModel  OvertimeSubtractionModel
	OvertimeSubtractionAmount float64
	HoursPerWeekdayException  WeekdayExceptionMap `gorm:"type:jsonb" sql:"json"`
	HolidayDaysPerYear        uint
	OvertimeWarningThreshold  uint
}

type WorkTimeModelCreateRequest struct {
	Name                      string                   `binding:"required"`
	WorkingHoursPerWeekday    float64                  `binding:"required"`
	OvertimeSubtractionModel  OvertimeSubtractionModel `binding:"required"`
	OvertimeSubtractionAmount float64
	HoursPerWeekdayException  WeekdayExceptionMap
	HolidayDaysPerYear        uint
	OvertimeWarningThreshold  uint
}

type WorkTimeModelUpdateRequest struct {
	WorkingHoursPerWeekday    float64
	OvertimeSubtractionModel  OvertimeSubtractionModel
	OvertimeSubtractionAmount float64
	HoursPerWeekdayException  WeekdayExceptionMap
	HolidayDaysPerYear        uint
	OvertimeWarningThreshold  uint
}
