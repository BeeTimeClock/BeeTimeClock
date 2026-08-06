package repository

import (
	"errors"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
)

type WorkTimeModel struct {
	env *core.Environment
}

func NewWorkTimeModel(env *core.Environment) *WorkTimeModel {
	return &WorkTimeModel{
		env: env,
	}
}

func (r *WorkTimeModel) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.WorkTimeModel{}, &model.UserWorktime{})
	if err != nil {
		return err
	}

	return nil
}

var ErrWorkTimeModelNotFound = errors.New("WorkTimeModel not found")

func (r WorkTimeModel) WorkTimeModelFindAll() ([]model.WorkTimeModel, error) {
	var items []model.WorkTimeModel
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

func (r WorkTimeModel) WorkTimeModelFindById(id uint) (model.WorkTimeModel, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.WorkTimeModel{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.WorkTimeModel
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.WorkTimeModel{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.WorkTimeModel{}, ErrWorkTimeModelNotFound
	}
	return item, result.Error
}

func (r WorkTimeModel) WorkTimeModelInsert(item *model.WorkTimeModel) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r WorkTimeModel) WorkTimeModelUpdate(item *model.WorkTimeModel) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r WorkTimeModel) WorkTimeModelDelete(item *model.WorkTimeModel) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

var ErrUserWorktimeNotFound = errors.New("UserWorktime not found")

func (r WorkTimeModel) UserWorktimeFindByUserID(userID uint) ([]model.UserWorktime, error) {
	var items []model.UserWorktime
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Preload("WorkTimeModel").Where("user_id = ?", userID).Order("valid_from DESC").Find(&items)
	if result.Error != nil {
		return items, result.Error
	}
	return items, nil
}

func (r WorkTimeModel) UserWorktimeFindByID(id uint) (model.UserWorktime, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.UserWorktime{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.UserWorktime
	result := db.Preload("WorkTimeModel").Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.UserWorktime{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.UserWorktime{}, ErrUserWorktimeNotFound
	}
	return item, nil
}

func (r WorkTimeModel) UserWorktimeInsert(item *model.UserWorktime) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r WorkTimeModel) UserWorktimeUpdate(item *model.UserWorktime) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Save(item)
	return result.Error
}

func (r WorkTimeModel) UserWorktimeDelete(item *model.UserWorktime) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}
