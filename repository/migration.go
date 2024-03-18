package repository

import (
	"errors"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
)

type Migration struct {
	env *core.Environment
}

func NewMigration(env *core.Environment) *Migration {
	return &Migration{
		env: env,
	}
}

func (r *Migration) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Migration{})
	if err != nil {
		return err
	}

	return nil
}

var ErrMigrationNotFound = errors.New("Migration not found")

func (r Migration) MigrationFindAll() ([]model.Migration, error) {
	var items []model.Migration
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

func (r Migration) MigrationFindById(id uint) (model.Migration, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Migration{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Migration
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Migration{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Migration{}, ErrMigrationNotFound
	}
	return item, result.Error
}

func (r Migration) MigrationFindByTitle(title string) (model.Migration, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Migration{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Migration
	result := db.Find(&item, "title = ?", title)
	if result.Error != nil {
		return model.Migration{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Migration{}, ErrMigrationNotFound
	}
	return item, result.Error
}

func (r Migration) MigrationInsert(item *model.Migration) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Migration) MigrationUpdate(item *model.Migration) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Migration) MigrationDelete(item *model.Migration) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}
