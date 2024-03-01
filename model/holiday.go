package model

import (
	"time"

	"gorm.io/gorm"
)

type Holiday struct {
	gorm.Model
	Name  string
	Date  time.Time `gorm:"uniqueIndex"`
	State string
}

type Holidays []Holiday

type HolidayImport struct {
	Datum   string `json:"datum"`
	Hinweis string `json:"hinweis"`
}

func (hi HolidayImport) GetDate() (time.Time, error) {
	return time.Parse("2006-01-02", hi.Datum)
}

func (h Holidays) Contains(date time.Time) bool {
	for _, holiday := range h {
		if holiday.Date == date {
			return true
		}
	}

	return false
}

func GetNeededHoursForMonth(holidays Holidays, year int, month int) float64 {
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	hours := 0.0

	currentDay := firstOfMonth

	for !currentDay.After(lastOfMonth) {
		skip := false

		if currentDay.Weekday() == time.Saturday || currentDay.Weekday() == time.Sunday {
			skip = true
		}

		if holidays.Contains(currentDay) {
			skip = true
		}

		if !skip {
			if worktime, exists := DefaultWorkTimeModel().HoursPerWeekdayException[currentDay.Weekday()]; exists {
				hours += worktime
			} else {
				hours += DefaultWorkTimeModel().DefaultHoursPerWeekday
			}
		}
		currentDay = currentDay.AddDate(0, 0, 1)
	}

	return hours
}
