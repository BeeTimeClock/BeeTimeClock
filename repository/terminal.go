package repository

import (
	"errors"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
)

type Terminal struct {
	env *core.Environment
}

func NewTerminal(env *core.Environment) *Terminal {
	return &Terminal{
		env: env,
	}
}

func (r Terminal) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Terminal{})
	if err != nil {
		return err
	}

	return nil
}

var ErrTerminalNotFound = errors.New("Terminal not found")

func (r Terminal) TerminalFindAll() ([]model.Terminal, error) {
	var items []model.Terminal
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r Terminal) TerminalFindById(id uint) (model.Terminal, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Terminal{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Terminal
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Terminal{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Terminal{}, ErrTerminalNotFound
	}
	return item, result.Error
}

func (r Terminal) TerminalInsert(item *model.Terminal) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Terminal) TerminalUpdate(item *model.Terminal) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Terminal) TerminalDelete(item *model.Terminal) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

func (r Terminal) TerminalFindByClientId(clientId string) (model.Terminal, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Terminal{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Terminal
	result := db.Find(&item, "client_id = ?", clientId)
	if result.Error != nil {
		return model.Terminal{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Terminal{}, ErrTerminalNotFound
	}
	return item, result.Error
}
