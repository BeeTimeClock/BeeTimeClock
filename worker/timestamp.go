package worker

import (
	"slices"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
)

type Timestamp struct {
	env          *core.Environment
	holiday      *repository.Holiday
	timestamp    *repository.Timestamp
	externalWork *repository.ExternalWork
	absence      *repository.Absence
	user         *repository.User
}

func NewTimestamp(env *core.Environment, user *repository.User, externalWork *repository.ExternalWork, timestamp *repository.Timestamp, holiday *repository.Holiday, absence *repository.Absence) *Timestamp {
	return &Timestamp{
		env:          env,
		holiday:      holiday,
		timestamp:    timestamp,
		externalWork: externalWork,
		user:         user,
		absence:      absence,
	}
}

func (w *Timestamp) CalculateMonth(userID uint, year int, month int) (model.TimestampMonthCalculated, error) {
	result := model.TimestampMonthCalculated{
		TimestampGroups: []model.TimestampGroup{},
	}
	currentLocation := time.Now().Location()

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	user, err := w.user.FindByID(userID)
	if err != nil {
		return result, err
	}
	timestamps, err := w.timestamp.FindByUserIDAndDate(userID, firstOfMonth, lastOfMonth)
	if err != nil {
		return result, err
	}

	holidays, err := w.holiday.HolidayFindByDateRange(firstOfMonth, lastOfMonth)
	if err != nil {
		return result, err
	}
	neededHours := model.GetNeededHoursForMonth(holidays, year, month)

	grouped := make(map[time.Time]model.TimestampGroup)

	for _, timestamp := range timestamps {
		year, month, day := timestamp.ComingTimestamp.Date()
		timestamp_date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		if _, exists := grouped[timestamp_date]; !exists {
			grouped[timestamp_date] = model.TimestampGroup{
				IsHomeoffice: true,
				Date:         timestamp_date,
			}
		}
		group := grouped[timestamp_date]

		if !timestamp.IsHomeoffice {
			group.IsHomeoffice = false
		}
		group.Timestamps = append(grouped[timestamp_date].Timestamps, timestamp)

		workingHours, subtractedHours := timestamp.CalculateWorkingHours()
		group.WorkingHours += workingHours
		group.SubtractedHours += subtractedHours

		workTimeModel := model.DefaultWorkTimeModel()

		neededHours := workTimeModel.GetWorkingHoursForDay(timestamp_date, holidays)

		group.OvertimeHours = group.WorkingHours - neededHours

		grouped[timestamp_date] = group
	}

	for _, value := range grouped {
		result.TimestampGroups = append(result.TimestampGroups, value)
		result.OvertimeHours += value.OvertimeHours
	}

	subtractedHours := 0.0
	if result.OvertimeHours > 0 {
		switch user.OvertimeSubtractionModel {
		case model.OVERTIME_SUBTRACTION_MODEL_HOURS:
			subtractedHours = user.OvertimeSubtractionAmount
			if result.OvertimeHours < user.OvertimeSubtractionAmount {
				subtractedHours = result.OvertimeHours
			}
			break
		case model.OVERTIME_SUBTRACTION_MODEL_PERCENTAGE:
			subtractionAmount := (neededHours / 100 * user.OvertimeSubtractionAmount)
			subtractedHours = subtractionAmount
			if result.OvertimeHours < subtractionAmount {
				subtractedHours = result.OvertimeHours
			}
			break
		}
	}

	result.SubtractedHours = subtractedHours

	return result, nil
}

func (w *Timestamp) MissingDays(userID uint) ([]time.Time, error) {
	user, err := w.user.FindByID(userID)
	if err != nil {
		return nil, err
	}
	lastTimestamp, err := w.timestamp.FindLastByUserID(userID)
	if err != nil {
		return nil, err
	}

	current := user.CreatedAt
	endDate := lastTimestamp.ComingTimestamp

	result := []time.Time{}
	for current.Before(endDate) || current.Equal(endDate) {
		entries, err := w.MissingDaysInMonth(userID, current.Year(), int(current.Month()))
		if err != nil {
			return nil, err
		}
		current = current.AddDate(0, 1, 0)
		result = append(result, entries...)
	}
	return result, nil
}

func (w *Timestamp) MissingDaysInMonth(userID uint, year int, month int) ([]time.Time, error) {
	holidays, err := w.holiday.HolidayFindByYear(year)
	if err != nil {
		return nil, err
	}

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-1 * time.Second)

	absences, err := w.absence.AbsenceFindByUserIDAndBetweenDates(userID, firstOfMonth, lastOfMonth)
	if err != nil {
		return nil, err
	}

	timestamps, err := w.timestamp.FindByUserIDAndDate(userID, firstOfMonth, lastOfMonth)
	if err != nil {
		return nil, err
	}

	externalWork, err := w.externalWork.ExternalWorkFindByUserIDAndEndBetween(userID, firstOfMonth, lastOfMonth)
	if err != nil {
		return nil, err
	}

	missingDays := []time.Time{}
	currentDay := firstOfMonth.AddDate(0, 0, -1)
	for currentDay.Before(lastOfMonth) {
		currentDay = currentDay.AddDate(0, 0, 1)
		if currentDay.After(time.Now()) {
			break
		}

		if currentDay.Weekday() == time.Sunday || currentDay.Weekday() == time.Saturday {
			continue
		}

		if holidays.Contains(currentDay) {
			continue
		}

		if slices.ContainsFunc(absences, func(n model.Absence) bool {
			return n.IsDateInAbsence(currentDay)
		}) {
			continue
		}

		if slices.ContainsFunc(externalWork, func(n model.ExternalWork) bool {
			return n.IsDateInExternalWork(currentDay)
		}) {
			continue
		}

		y1, m1, d1 := currentDay.Date()
		if slices.ContainsFunc(timestamps, func(n model.Timestamp) bool {
			y2, m2, d2 := n.ComingTimestamp.Date()

			return y1 == y2 && m1 == m2 && d1 == d2
		}) {
			continue
		}

		missingDays = append(missingDays, currentDay)
	}

	return missingDays, nil
}
