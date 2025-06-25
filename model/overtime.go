package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const (
	OVERTIME_SUMMARY_ENTRY_SOURCE_TIMESTAMP     OvertimeSummaryEntrySource = "timestamp"
	OVERTIME_SUMMARY_ENTRY_SOURCE_EXTERNAL_WORK OvertimeSummaryEntrySource = "external_work"
	OVERTIME_SUMMARY_ENTRY_SOURCE_ABSENCE       OvertimeSummaryEntrySource = "absence"
)

type OvertimeSummaryEntrySource string

type OvertimeSummaryEntry struct {
	Source     string
	Identifier *uint
	Value      float64
}

type OvertimeSummary []OvertimeSummaryEntry

type OvertimeMonthQuota struct {
	gorm.Model
	UserID  uint `gorm:"index:idx_month_quota,unique;index"`
	User    User
	Year    int `gorm:"index:idx_month_quota,unique"`
	Month   int `gorm:"index:idx_month_quota,unique"`
	Hours   *float64
	Summary OvertimeSummary `gorm:"type:jsonb" sql:"json"`
}

func (o *OvertimeMonthQuota) InsertSummary(source string, identifier *uint, value float64) {
	o.Summary = append(o.Summary, OvertimeSummaryEntry{
		Source:     source,
		Identifier: identifier,
		Value:      value,
	})
}

func (o *OvertimeMonthQuota) Calculate() {
	hours := 0.0
	for _, entry := range o.Summary {
		hours += entry.Value
	}
	o.Hours = &hours
}

func (o *OvertimeSummary) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := &OvertimeSummary{}
	err := json.Unmarshal(bytes, &result)
	*o = OvertimeSummary(*result)
	return err
}

func (o OvertimeSummary) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
