package repository

import (
	"fmt"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
)

type Fuel struct {
	env *core.Environment
}

func NewFuel(env *core.Environment) *Fuel {
	return &Fuel{
		env: env,
	}
}

func (r *Fuel) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Fuel{})
	if err != nil {
		return err
	}

	return nil
}

func (r *Fuel) FindAll() ([]model.Fuel, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Fuel

	result := db.Find(&items)
	return items, result.Error
}

func (r *Fuel) FindByID(id uint) (model.Fuel, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Fuel{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Fuel
	result := db.Find(&item, "id = ?", id)

	if result.RowsAffected == 0 {
		return model.Fuel{}, fmt.Errorf("no fuel with id %d found", id)
	}

	return item, result.Error
}

func (r *Fuel) FindByUserID(userID uint) ([]model.Fuel, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Fuel
	result := db.Find(&items, "user_id = ?", userID)

	return items, result.Error
}

func (r *Fuel) FindByUserIDAndState(userID uint, state model.FuelState) ([]model.Fuel, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Fuel
	result := db.Find(&items, "user_id = ? and state = ?", userID, state)

	return items, result.Error
}

func (r *Fuel) FindLastByUserID(userID uint) (model.Fuel, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Fuel{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Fuel
	result := db.Last(&item, "user_id = ?", userID)

	return item, result.Error
}

func (r *Fuel) FindByUserIDAndDate(userID uint, from, till time.Time) ([]model.Fuel, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Fuel
	result := db.Find(&items, "user_id = ? AND coming_fuel BETWEEN ? AND ?", userID, from, till)

	return items, result.Error
}

func (r *Fuel) CountByUserID(userID uint) (int64, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return 0, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var count int64
	result := db.Model(&model.Fuel{}).Where("user_id = ?", userID).Count(&count)
	return count, result.Error
}

func (r *Fuel) Insert(fuel *model.Fuel) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(&fuel)
	return result.Error
}

func (r *Fuel) Update(fuel *model.Fuel) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(&fuel)
	return result.Error
}

func (r *Fuel) Delete(fuel *model.Fuel) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(&fuel)
	return result.Error
}
