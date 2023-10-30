package repository

import (
	"fmt"

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

	err = db.AutoMigrate(&model.AbsenceReason{})
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
	result = result.Where(query, args).Find(&items)
	return items, result.Error
}

func (r *Absence) FindByID(id uint) (model.Absence, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Absence{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Absence
	result := db.Find(&item, "id = ?", id)

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
