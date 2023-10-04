package core

import (
	"github.com/BeeTimeClock/BeeTimeClock-Server/database"
)

type Environment struct {
	DatabaseManager *database.DatabaseManager
	UploadPath      string
	Secret          []byte
}

func NewEnvironment() *Environment {
	return &Environment{
		UploadPath: "upload",
	}
}
