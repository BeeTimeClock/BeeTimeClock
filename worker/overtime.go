package worker

import (
	"slices"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
)

type Overtime struct {
	env             *core.Environment
	holiday         *repository.Holiday
	timestamp       *repository.Timestamp
	externalWork    *repository.ExternalWork
	overtime        *repository.Overtime
	user            *repository.User
	absence         *repository.Absence
	timestampWorker *Timestamp
}

func NewOvertime(env *core.Environment, user *repository.User, externalWork *repository.ExternalWork, timestamp *repository.Timestamp, holiday *repository.Holiday, overtime *repository.Overtime, timestampWorker *Timestamp, absence *repository.Absence) *Overtime {
	return &Overtime{
		env:             env,
		holiday:         holiday,
		timestamp:       timestamp,
		externalWork:    externalWork,
		user:            user,
		overtime:        overtime,
		timestampWorker: timestampWorker,
		absence:         absence,
	}
}

func (w *Overtime) CalculateMonth(userID uint, year int, month int) (model.OvertimeMonthQuota, bool, error) {
	hours := 0.0
	result := model.OvertimeMonthQuota{
		UserID:  userID,
		Year:    year,
		Month:   month,
		Hours:   &hours,
		Summary: model.OvertimeSummary{},
	}

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-1 * time.Second)
	holidays, err := w.holiday.HolidayFindByDateRange(firstOfMonth, lastOfMonth)
	if err != nil {
		return model.OvertimeMonthQuota{}, false, err
	}

	timestampOvertime, err := w.timestampWorker.CalculateMonth(userID, year, month)
	if err != nil {
		return model.OvertimeMonthQuota{}, false, err
	}
	result.InsertSummary("timestamp_raw", nil, timestampOvertime.OvertimeHours, 1.0)
	result.InsertSummary("timestamp_subtracted", nil, timestampOvertime.SubtractedHours*-1.0, 1.0)
	result.InsertSummary("timestamp_final", nil, timestampOvertime.OvertimeHours-timestampOvertime.SubtractedHours, 0)

	externalWorks, err := w.externalWork.ExternalWorkFindByUserIDAndEndBetween(userID, firstOfMonth, lastOfMonth)
	if err != nil {
		return model.OvertimeMonthQuota{}, false, err
	}

	for _, externalWork := range externalWorks {
		calculated := externalWork.Calculate(holidays)
		result.InsertSummary("external_work", &externalWork.ID, calculated.TotalOvertimeHours, 1.0)
	}

	absences, err := w.absence.AbsenceFindByUserIDAndBetweenDates(userID, firstOfMonth, lastOfMonth)
	if err != nil {
		return model.OvertimeMonthQuota{}, false, err
	}
	for _, absence := range absences {
		switch absence.AbsenceReason.Impact {
		case model.ABESENCE_REASON_OVERTIME_IMPACT_DURATION:
			duration := absence.AbsenceTill.Sub(absence.AbsenceFrom).Hours()

			result.InsertSummary("absence", &absence.ID, duration, 1.0)
			break
		}
	}

	result.Calculate()

	quota, err := w.overtime.OvertimeMonthQuotaFindByUserIDAndYearAndMonth(userID, year, month)
	quotaNotExists := err == repository.ErrOvertimeMonthQuotaNotFound

	if !quotaNotExists && err != nil {
		return model.OvertimeMonthQuota{}, false, err
	}

	if quotaNotExists {
		err = w.overtime.OvertimeMonthQuotaInsert(&result)
		if err != nil {
			return model.OvertimeMonthQuota{}, false, err
		}

		return result, true, nil
	} else {
		quota.Hours = result.Hours
		quota.Summary = result.Summary

		err = w.overtime.OvertimeMonthQuotaUpdate(&quota)
		if err != nil {
			return model.OvertimeMonthQuota{}, false, err
		}
	}

	return result, false, nil
}

func (w *Overtime) CalculateMissingMonths() error {
	timestamps, err := w.timestamp.FindYearMonthsWithTimestamps()
	if err != nil {
		return err
	}
	overtimeMonths, err := w.overtime.OvertimeMonthQuotaFindAll()
	if err != nil {
		return err
	}

	for _, timestamp := range timestamps {
		exists := slices.ContainsFunc(overtimeMonths, func(n model.OvertimeMonthQuota) bool {
			return n.UserID == timestamp.UserID && n.Year == timestamp.Year && n.Month == timestamp.Month
		})

		if !exists {
			w.CalculateMonth(timestamp.UserID, timestamp.Year, timestamp.Month)
		}
	}

	return nil
}
