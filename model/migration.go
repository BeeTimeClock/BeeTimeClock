package model

import (
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	gorm.Model

	Title      string
	Result     string
	FinishedAt time.Time
	Success    bool
}
