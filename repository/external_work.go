package repository

import (
	"errors"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"gorm.io/gorm/clause"
)

type ExternalWork struct {
	env *core.Environment
}

func NewExternalWork(env *core.Environment) *ExternalWork {
	return &ExternalWork{
		env: env,
	}
}

func (r *ExternalWork) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.ExternalWork{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.ExternalWorkExpense{})
	if err != nil {
		return err
	}

	return nil
}

var ErrExternalWorkNotFound = errors.New("ExternalWork not found")

func (r ExternalWork) ExternalWorkFindAll() ([]model.ExternalWork, error) {
	var items []model.ExternalWork
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

func (r ExternalWork) ExternalWorkFindById(id uint, with_associations bool) (model.ExternalWork, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.ExternalWork{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.ExternalWork

	if with_associations {
		db = db.Preload(clause.Associations)
	}

	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.ExternalWork{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.ExternalWork{}, ErrExternalWorkNotFound
	}
	return item, result.Error
}

func (r ExternalWork) ExternalWorkFindByUserID(userId uint) ([]model.ExternalWork, error) {
	var items []model.ExternalWork
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items, "user_id = ?", userId)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r ExternalWork) ExternalWorkInsert(item *model.ExternalWork) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkUpdate(item *model.ExternalWork) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkDelete(item *model.ExternalWork) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

var ErrExternalWorkExpenseNotFound = errors.New("ExternalWorkExpense not found")

func (r ExternalWork) ExternalWorkExpenseFindAll() ([]model.ExternalWorkExpense, error) {
	var items []model.ExternalWorkExpense
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

func (r ExternalWork) ExternalWorkExpenseFindById(id uint) (model.ExternalWorkExpense, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.ExternalWorkExpense{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.ExternalWorkExpense
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.ExternalWorkExpense{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.ExternalWorkExpense{}, ErrExternalWorkExpenseNotFound
	}
	return item, result.Error
}

func (r ExternalWork) ExternalWorkExpenseFindByExternalWorkId(externalWorkId uint) ([]model.ExternalWorkExpense, error) {
	var items []model.ExternalWorkExpense
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items, "external_work_id = ?", externalWorkId)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r ExternalWork) ExternalWorkExpenseInsert(item *model.ExternalWorkExpense) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkExpenseUpdate(item *model.ExternalWorkExpense) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkExpenseDelete(item *model.ExternalWorkExpense) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

var ErrExternalWorkExternalEventNotFound = errors.New("ExternalWorkExternalEvent not found")

func (r ExternalWork) ExternalWorkExternalEventFindAll() ([]model.ExternalWorkExternalEvent, error) {
	var items []model.ExternalWorkExternalEvent
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

func (r ExternalWork) ExternalWorkExternalEventFindById(id uint) (model.ExternalWorkExternalEvent, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.ExternalWorkExternalEvent{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.ExternalWorkExternalEvent
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.ExternalWorkExternalEvent{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.ExternalWorkExternalEvent{}, ErrExternalWorkExternalEventNotFound
	}
	return item, result.Error
}

func (r ExternalWork) ExternalWorkExternalEventInsert(item *model.ExternalWorkExternalEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkExternalEventUpdate(item *model.ExternalWorkExternalEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkExternalEventDelete(item *model.ExternalWorkExternalEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}
