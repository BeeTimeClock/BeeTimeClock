package repository

import (
	"errors"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
)

type Holiday struct {
	env *core.Environment
}

func NewHoliday(env *core.Environment) *Holiday {
	return &Holiday{
		env: env,
	}
}

func (r *Holiday) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Holiday{}, &model.HolidayCustom{})
	if err != nil {
		return err
	}

	return nil
}

var ErrHolidayNotFound = errors.New("Holiday not found")

func (r Holiday) HolidayFindAll() ([]model.Holiday, error) {
	var items []model.Holiday
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

func (r Holiday) HolidayFindById(id uint) (model.Holiday, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Holiday{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Holiday
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.Holiday{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Holiday{}, ErrHolidayNotFound
	}
	return item, result.Error
}

func (r Holiday) HolidayInsert(item *model.Holiday) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Holiday) HolidayUpdate(item *model.Holiday) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Holiday) HolidayDelete(item *model.Holiday) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

func (r Holiday) HolidayFindByDate(date time.Time) (model.Holiday, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Holiday{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Holiday
	result := db.Find(&item, "date = ?", date)
	if result.Error != nil {
		return model.Holiday{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.Holiday{}, ErrHolidayNotFound
	}
	return item, result.Error
}

func (r Holiday) HolidayIsByDate(date time.Time) (bool, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return false, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Holiday
	result := db.Find(&item, "date = ?", date)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, result.Error
}

func (r Holiday) HolidayFindByYear(year int) (model.Holidays, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items model.Holidays
	result := db.Find(&items, "extract(year from date) = ?", year)
	if result.Error != nil {
		return nil, result.Error
	}

	return items, result.Error
}

func (r Holiday) HolidayFindByDateRange(start time.Time, end time.Time) ([]model.Holiday, error) {
	var items []model.Holiday
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items, "date between ? and ?", start, end)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

var ErrHolidayCustomNotFound = errors.New("HolidayCustom not found")

func (r Holiday) HolidayCustomFindAll() ([]model.HolidayCustom, error) {
	var items []model.HolidayCustom
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

func (r Holiday) HolidayCustomFindById(id uint) (model.HolidayCustom, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.HolidayCustom{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.HolidayCustom
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.HolidayCustom{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.HolidayCustom{}, ErrHolidayCustomNotFound
	}
	return item, result.Error
}

func (r Holiday) HolidayCustomInsert(item *model.HolidayCustom) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Holiday) HolidayCustomUpdate(item *model.HolidayCustom) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Holiday) HolidayCustomDelete(item *model.HolidayCustom) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}
