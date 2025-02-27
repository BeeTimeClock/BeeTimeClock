package core

import (
	"os"

	"github.com/BeeTimeClock/BeeTimeClock-Server/database"
)

type EnvironmentNotification struct {
	Enabled    bool
	WebhookUrl string
}

type EnvironmentStorage struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
}

func (es *EnvironmentStorage) HasS3() bool {
	return es.Endpoint != ""
}

type Environment struct {
	DatabaseManager *database.DatabaseManager
	UploadPath      string
	Secret          []byte
	Notification    EnvironmentNotification
	Storage         EnvironmentStorage
}

func NewEnvironment() *Environment {
	return &Environment{
		UploadPath: "upload",
		Notification: EnvironmentNotification{
			Enabled:    os.Getenv("NOTIFY_WEBHOOK_URL") != "",
			WebhookUrl: os.Getenv("NOTIFY_WEBHOOK_URL"),
		},
		Storage: EnvironmentStorage{
			Endpoint:        os.Getenv("BUCKET_ADDRESS"),
			AccessKeyID:     os.Getenv("BUCKET_USER"),
			SecretAccessKey: os.Getenv("BUCKET_PASSWORD"),
			BucketName:      os.Getenv("BUCKET_NAME"),
		},
	}
}
