package repository

import (
	"errors"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
)

type Overtime struct {
	env *core.Environment
}

func NewOvertime(env *core.Environment) *Overtime {
	return &Overtime{
		env: env,
	}
}

func (r *Overtime) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.OvertimeMonthQuota{})
	if err != nil {
		return err
	}

	return nil
}

var ErrOvertimeMonthQuotaNotFound = errors.New("OvertimeMonthQuota not found")

func (r Overtime) OvertimeMonthQuotaFindAll() ([]model.OvertimeMonthQuota, error) {
	var items []model.OvertimeMonthQuota
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

func (r Overtime) OvertimeMonthQuotaFindById(id uint) (model.OvertimeMonthQuota, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.OvertimeMonthQuota{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.OvertimeMonthQuota
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.OvertimeMonthQuota{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.OvertimeMonthQuota{}, ErrOvertimeMonthQuotaNotFound
	}
	return item, result.Error
}

func (r Overtime) OvertimeMonthQuotaInsert(item *model.OvertimeMonthQuota) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r Overtime) OvertimeMonthQuotaUpdate(item *model.OvertimeMonthQuota) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r Overtime) OvertimeMonthQuotaDelete(item *model.OvertimeMonthQuota) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

func (r Overtime) OvertimeMonthQuotaFindByUserID(userID uint) ([]model.OvertimeMonthQuota, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.OvertimeMonthQuota
	result := db.Find(&items, "user_id = ?", userID)

	return items, result.Error
}

func (r Overtime) OvertimeMonthQuotaFindByUserIDAndYearAndMonth(userID uint, year int, month int) (model.OvertimeMonthQuota, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.OvertimeMonthQuota{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.OvertimeMonthQuota
	result := db.Find(&item, "user_id = ? and year = ? and month = ?", userID, year, month)
	if result.Error != nil {
		return model.OvertimeMonthQuota{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.OvertimeMonthQuota{}, ErrOvertimeMonthQuotaNotFound
	}
	return item, result.Error
}
