package repository

import (
	"errors"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/google/uuid"
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

	err = db.AutoMigrate(&model.ExternalWorkCompensation{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.ExternalWorkExternalEvent{})
	if err != nil {
		return err
	}

	externalWorkCompensation := model.ExternalWorkCompensation{
		IsoCountryCodeA2:         "DE",
		PrivateCarKmCompensation: 0.3,
		WithSocialInsuranceSlots: []model.ExternalWorkCompensationHourSlot{
			{
				Hours:        6,
				Compensation: 8.00,
			},
			{
				Hours:        8,
				Compensation: 13.00,
			},
			{
				Hours:        10,
				Compensation: 15.00,
			},
			{
				Hours:        14,
				Compensation: 16.00,
			},
			{
				Hours:        24,
				Compensation: 0.0,
			},
		},
		WithoutSocialInsuranceSlots: []model.ExternalWorkCompensationHourSlot{
			{
				Hours:        8,
				Compensation: 14.00,
			},
			{
				Hours:        24,
				Compensation: 28.00,
			},
		},
		AdditionalOptions: map[string]float64{
			"Aussendienstpauschale": 20.00,
		},
		ValidFrom: time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.Local),
	}

	_, err = r.ExternalWorkCompensationFindByCountryCode("DE")
	if err != nil {
		if err == ErrExternalWorkCompensationNotFound {
			r.ExternalWorkCompensationInsert(&externalWorkCompensation)
		} else {
			return err
		}
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

func (r ExternalWork) ExternalWorkFindByUserIDAndInvoiceIdentifier(userId uint, invoiceIdentifier uuid.UUID) ([]model.ExternalWork, error) {
	var items []model.ExternalWork
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items, "user_id = ? and invoice_identifier = ?", userId, invoiceIdentifier)
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

	////	DepartureTime       *time.Time
	////	ArrivalTime         *time.Time
	////	TravelDurationHours float64
	////	PauseDurationHours  float64
	////	RestDurationHours   float64
	////	OnSiteFrom          *time.Time
	////	OnSiteTill          *time.Time
	//	Place               string
	payload := map[string]any{
		"departure_time":        item.DepartureTime,
		"arrival_time":          item.ArrivalTime,
		"travel_duration_hours": item.TravelDurationHours,
		"pause_duration_hours":  item.PauseDurationHours,
		"on_site_from":          item.OnSiteFrom,
		"on_site_till":          item.OnSiteTill,
		"place":                 item.Place,
		"additional_options":    item.AdditionalOptions,
	}

	result := db.Debug().Model(item).Updates(payload)
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

func (r ExternalWork) ExternalWorkFindByUserIDAndEndBetween(userId uint, start time.Time, end time.Time) ([]model.ExternalWork, error) {
	var items []model.ExternalWork
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Preload(clause.Associations).Find(&items, "user_id = ? and till between ? and ?", userId, start, end)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r ExternalWork) ExternalWorkFindByUserIDAndStatus(userId uint, status model.ExternalWorkStatus) ([]model.ExternalWork, error) {
	var items []model.ExternalWork
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Preload(clause.Associations).Find(&items, "user_id = ? and status = ?", userId, status)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

var ErrExternalWorkCompensationNotFound = errors.New("ExternalWorkCompensation not found")

func (r ExternalWork) ExternalWorkCompensationFindAll() ([]model.ExternalWorkCompensation, error) {
	var items []model.ExternalWorkCompensation
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Preload(clause.Associations).Find(&items)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r ExternalWork) ExternalWorkCompensationFindById(id uint) (model.ExternalWorkCompensation, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.ExternalWorkCompensation{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.ExternalWorkCompensation
	result := db.Preload(clause.Associations).Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.ExternalWorkCompensation{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.ExternalWorkCompensation{}, ErrExternalWorkCompensationNotFound
	}
	return item, result.Error
}

func (r ExternalWork) ExternalWorkCompensationInsert(item *model.ExternalWorkCompensation) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkCompensationUpdate(item *model.ExternalWorkCompensation) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkCompensationDelete(item *model.ExternalWorkCompensation) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(item)
	return result.Error
}

func (r ExternalWork) ExternalWorkCompensationFindByCountryCode(countryCode string) (model.ExternalWorkCompensation, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.ExternalWorkCompensation{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.ExternalWorkCompensation
	result := db.Find(&item, "iso_country_code_a2 = ?", countryCode)
	if result.Error != nil {
		return model.ExternalWorkCompensation{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.ExternalWorkCompensation{}, ErrExternalWorkCompensationNotFound
	}
	return item, result.Error
}
