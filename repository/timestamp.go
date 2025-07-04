package repository

import (
	"errors"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"gorm.io/gorm/clause"
)

type Timestamp struct {
	env *core.Environment
}

var ErrTimestampNotFound = errors.New("timestamp not found")

func NewTimestamp(env *core.Environment) *Timestamp {
	return &Timestamp{
		env: env,
	}
}

func (r *Timestamp) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.Timestamp{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.TimestampCorrection{})
	if err != nil {
		return err
	}

	return nil
}

func (r *Timestamp) FindAll() ([]model.Timestamp, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Timestamp

	result := db.Find(&items)
	return items, result.Error
}

func (r *Timestamp) FindByID(id uint) (model.Timestamp, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Timestamp{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Timestamp
	result := db.Find(&item, "id = ?", id)

	if result.RowsAffected == 0 {
		return model.Timestamp{}, ErrTimestampNotFound
	}

	return item, result.Error
}

func (r *Timestamp) FindByUserID(userID uint) ([]model.Timestamp, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Timestamp
	result := db.Find(&items, "user_id = ?", userID)

	return items, result.Error
}

func (r *Timestamp) FindLastByUserID(userID uint) (model.Timestamp, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.Timestamp{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.Timestamp
	result := db.Order("coming_timestamp DESC").Last(&item, "user_id = ?", userID)

	return item, result.Error
}

func (r *Timestamp) FindSuspiciousTimestampsByUserID(userId uint) ([]model.Timestamp, error) {
	var items []model.Timestamp
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Preload(clause.Associations).Find(&items, "user_id = ? and (date_part('year', coming_timestamp) < 1999 or date_part('year', going_timestamp) < 1999)", userId)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}

func (r *Timestamp) FindByUserIDAndDate(userID uint, from, till time.Time) ([]model.Timestamp, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.Timestamp
	result := db.Preload(clause.Associations).Order("coming_timestamp DESC").Find(&items, "user_id = ? AND coming_timestamp BETWEEN ? AND ?", userID, from, till)

	return items, result.Error
}

func (r *Timestamp) CountByUserID(userID uint) (int64, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return 0, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var count int64
	result := db.Model(&model.Timestamp{}).Where("user_id = ?", userID).Count(&count)
	return count, result.Error
}

func (r *Timestamp) Insert(timestamp *model.Timestamp) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(&timestamp)
	return result.Error
}

func (r *Timestamp) Update(timestamp *model.Timestamp) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(&timestamp)
	return result.Error
}

func (r *Timestamp) Delete(timestamp *model.Timestamp) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(&timestamp)
	return result.Error
}

func (r *Timestamp) TimestampCorrectionInsert(timestampCorrection *model.TimestampCorrection) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(&timestampCorrection)
	return result.Error
}

func (r *Timestamp) TimestampCorrectionFindByTimestampID(timestampID uint) ([]model.TimestampCorrection, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.TimestampCorrection
	result := db.Find(&items, "timestamp_id = ?", timestampID)

	return items, result.Error
}

func (r *Timestamp) FindYearMonthsWithTimestampsByUserId(userID uint) ([]model.TimestampYearMonthGrouped, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.TimestampYearMonthGrouped

	var count int64
	result := db.Model(&model.Timestamp{}).Where("user_id = ?", userID).Count(&count)

	if count == 0 {
		return items, result.Error
	}

	result = db.Model(&model.Timestamp{}).
		Select("distinct extract(year from coming_timestamp) as year, extract(month from coming_timestamp) as month").
		Where("user_id = ?", userID).
		Scan(&items)

	return items, result.Error
}

func (r *Timestamp) FindYearMonthsWithTimestamps() ([]model.TimestampYearMonthGrouped, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.TimestampYearMonthGrouped

	var count int64
	result := db.Model(&model.Timestamp{}).Count(&count)

	if count == 0 {
		return items, result.Error
	}

	result = db.Model(&model.Timestamp{}).
		Select("distinct extract(year from coming_timestamp) as year, extract(month from coming_timestamp) as month, user_id").
		Scan(&items)

	return items, result.Error
}
