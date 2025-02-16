package core

import (
	"os"

	"github.com/BeeTimeClock/BeeTimeClock-Server/database"
)

type EnvironmentNotification struct {
	Enabled    bool
	WebhookUrl string
}

type Environment struct {
	DatabaseManager *database.DatabaseManager
	UploadPath      string
	Secret          []byte
	Notification    EnvironmentNotification
}

func NewEnvironment() *Environment {
	return &Environment{
		UploadPath: "upload",
		Notification: EnvironmentNotification{
			Enabled:    os.Getenv("NOTIFY_WEBHOOK_URL") != "",
			WebhookUrl: os.Getenv("NOTIFY_WEBHOOK_URL"),
		},
	}
}
