package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"gorm.io/gorm/clause"
)

type Absence struct {
	env *core.Environment
}

func NewAbsence(env *core.Environment) *Absence {
	return &Absence{
		env: env,
	}
}

func (r *Absence) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Absence{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.AbsenceExternalEvent{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.AbsenceReason{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.Holiday{})
	if err != nil {
		return err
	}

	absenceReasons := []string{
		"Krank (mit AU)",
		"Krank (ohne AU)",
		"Sonderurlaub",
		"Urlaub",
		"Berufsschule",
		"Aussendienst",
	}

	existingReasons, err := r.FindAllAbsenceReasons()
	if err != nil {
		return err
	}

	for _, reason := range absenceReasons {
		reasonExists := false
		for _, existingReason := range existingReasons {
			if existingReason.Description == reason {
				reasonExists = true
				break
			}
		}

		if !reasonExists {
			err = r.InsertAbsenceReason(&model.AbsenceReason{
				Description: reason,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Absence) FindAll(withUser bool) ([]model.Absence, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Absence

	result := db
	if withUser {
		result = result.Preload(clause.Associations)
	}
	result = result.Find(&items)
	return items, result.Error
}

func (r *Absence) FindByQuery(withUser bool, query string, args ...interface{}) ([]model.Absence, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Absence

	result := db
	if withUser {
		result = result.Preload(clause.Associations)
	}
	result = result.Where(query, args...).Find(&items)
	return items, result.Error
}

func (r *Absence) FindByID(id uint) (model.Absence, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Absence{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Absence
	result := db.Preload(clause.Associations).Find(&item, "id = ?", id)

	if result.RowsAffected == 0 {
		return model.Absence{}, fmt.Errorf("no absence with id %d found", id)
	}

	return item, result.Error
}

func (r *Absence) FindByUserID(userID uint) ([]model.Absence, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Absence
	result := db.Find(&items, "user_id = ?", userID)

	return items, result.Error
}

func (r *Absence) FindByUserIDAndYear(userID uint, year int) ([]model.Absence, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Absence
	result := db.Find(&items, "user_id = ? and absence_from between ? and ?", userID,
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(year, 12, 31, 23, 59, 0, 0, time.UTC))

	return items, result.Error
}

func (r *Absence) Insert(absence *model.Absence) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(&absence)
	return result.Error
}

func (r *Absence) Update(absence *model.Absence) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(&absence)
	return result.Error
}

func (r *Absence) Delete(absence *model.Absence) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(&absence)
	return result.Error
}

func (r *Absence) FindAllAbsenceReasons() ([]model.AbsenceReason, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.AbsenceReason
	result := db.Find(&items)

	return items, result.Error
}

func (r *Absence) InsertAbsenceReason(absenceReason *model.AbsenceReason) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(absenceReason)
	return result.Error
}

func (r *Absence) FindAbsenceReasonByID(id uint) (model.AbsenceReason, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.AbsenceReason{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.AbsenceReason
	result := db.Find(&item, "id = ?", id)

	if result.RowsAffected == 0 {
		return model.AbsenceReason{}, fmt.Errorf("no absence with id %d found", id)
	}

	return item, result.Error
}

var ErrHolidayNotFound = errors.New("Holiday not found")

func (r Absence) HolidayFindAll() ([]model.Holiday, error) {
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

func (r Absence) HolidayFindById(id uint) (model.Holiday, error) {
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

func (r Absence) HolidayInsert(item *model.Holiday) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Absence) HolidayUpdate(item *model.Holiday) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Absence) HolidayDelete(item *model.Holiday) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

func (r Absence) HolidayFindByDate(date time.Time) (model.Holiday, error) {
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

func (r Absence) HolidayIsByDate(date time.Time) (bool, error) {
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

func (r Absence) HolidayFindByYear(year int) (model.Holidays, error) {
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

func (r Absence) FindYearsWithAbsencesByUserId(userID uint) ([]int, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []int

	result := db.Model(&model.Absence{}).
		Select("distinct extract(year from absence_from) as year").
		Where("user_id = ?", userID).
		Scan(&items)

	return items, result.Error
}

var ErrAbsenceExternalEventNotFound = errors.New("AbsenceExternalEvent not found")

func (r Absence) AbsenceExternalEventFindAll() ([]model.AbsenceExternalEvent, error) {
	var items []model.AbsenceExternalEvent
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

func (r Absence) AbsenceExternalEventFindById(id uint) (model.AbsenceExternalEvent, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.AbsenceExternalEvent{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.AbsenceExternalEvent
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.AbsenceExternalEvent{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.AbsenceExternalEvent{}, ErrAbsenceExternalEventNotFound
	}
	return item, result.Error
}

func (r Absence) AbsenceExternalEventInsert(item *model.AbsenceExternalEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Absence) AbsenceExternalEventUpdate(item *model.AbsenceExternalEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Absence) AbsenceExternalEventDelete(item *model.AbsenceExternalEvent) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(item)
	return result.Error
}

func (r Absence) AbsenceExternalEventFindByAbsenceId(absenceId uint) ([]model.AbsenceExternalEvent, error) {
	var items []model.AbsenceExternalEvent
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items, "absence_id = ?", absenceId)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}
