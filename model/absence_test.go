package model

import (
	"testing"
	"time"
)

func TestAbsenceDaysCalculation(t *testing.T) {
	type AbsenceTestData struct {
		From time.Time
		Till time.Time
		Wanted int
	}


	testData := []AbsenceTestData{
		{
			From: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			Till: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			Wanted: 1,
		},
		{
			From: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			Till: time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
			Wanted: 8,
		},
		{
			From: time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
			Till: time.Date(2024, 1, 12, 0, 0, 0, 0, time.UTC),
			Wanted: 5,
		},
		{
			From: time.Date(2024, 07, 15, 0, 0, 0, 0, time.UTC),
			Till: time.Date(2024, 07, 30, 0, 0, 0, 0, time.UTC),
			Wanted: 12,
		},
	}

	
	for _, item := range testData {
		absence := Absence{
			AbsenceFrom: item.From,
			AbsenceTill: item.Till,
		}

		workdays := absence.GetAbsenceWorkDays()
		if workdays != item.Wanted {
			t.Fatalf("From: %s, Till: %s, Want: %d, Got: %d\n", item.From, item.Till, item.Wanted, workdays)
		} else {
			t.Logf("From: %s, Till: %s, Want: %d, Got: %d\n", item.From, item.Till, item.Wanted, workdays)
		}
	}
}
