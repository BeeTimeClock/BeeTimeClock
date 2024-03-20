package repository

import (
	"errors"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"gorm.io/gorm/clause"
)

type Settings struct {
	env *core.Environment
}

func NewSettings(env *core.Environment) *Settings {
	return &Settings{
		env: env,
	}
}

func (r Settings) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Settings{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.SettingsOfficeIPAddresses{})
	if err != nil {
		return err
	}

	_, err = r.SettingsFind()
	if err != nil {
		return err
	}

	return nil
}

var ErrSettingsNotFound = errors.New("Settings not found")

func (r Settings) SettingsFind() (model.Settings, error) {
	var item model.Settings
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return item, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Preload(clause.Associations).FirstOrCreate(&item)
	if result.Error != nil {
		return item, result.Error
	}
	return item, result.Error
}

func (r Settings) SettingsUpdate(item *model.Settings) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}
