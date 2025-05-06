package model

import (
	"testing"
	"time"
)

func TestExternalWorkExpense_Calculate(t *testing.T) {
	departure := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
	arriaval := time.Date(2025, 1, 1, 18, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string // description of this test case
		input ExternalWorkExpense
		want  ExternalWorkExpenseCalculated
	}{
		{
			name: "No Arriaval on day",
			input: ExternalWorkExpense{
				Date:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				DepartureTime: &departure,
			},
			want: ExternalWorkExpenseCalculated{
				TotalAwayHours: 15,
			},
		},
		{
			name: "No Departure on day",
			input: ExternalWorkExpense{
				Date:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				ArrivalTime: &arriaval,
			},
			want: ExternalWorkExpenseCalculated{
				TotalAwayHours: 18,
			},
		},
		{
			name: "Departure same day",
			input: ExternalWorkExpense{
				Date:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				DepartureTime: &departure,
				ArrivalTime:   &arriaval,
			},
			want: ExternalWorkExpenseCalculated{
				TotalAwayHours: 9,
			},
		},
		{
			name: "Complete Day",
			input: ExternalWorkExpense{
				Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: ExternalWorkExpenseCalculated{
				TotalAwayHours: 24,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Calculate()
			if got.TotalAwayHours != tt.want.TotalAwayHours {
				t.Errorf("TotalAwayHours = %v, want %v", got.TotalAwayHours, tt.want.TotalAwayHours)
			}
		})
	}
}
