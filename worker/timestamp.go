package worker

import (
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
	user         *repository.User
}

func NewTimestamp(env *core.Environment, user *repository.User, externalWork *repository.ExternalWork, timestamp *repository.Timestamp, holiday *repository.Holiday) *Timestamp {
	return &Timestamp{
		env:          env,
		holiday:      holiday,
		timestamp:    timestamp,
		externalWork: externalWork,
		user:         user,
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
