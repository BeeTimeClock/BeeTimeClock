package repository

import (
	"errors"
	"fmt"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"gorm.io/gorm/clause"
)

type User struct {
	env *core.Environment
}

var ErrUserNotFound = errors.New("user not found")
var ErrUserApikeyNotFound = errors.New("user apikey not found")

func NewUser(env *core.Environment) *User {
	return &User{
		env: env,
	}
}

func (r *User) Migrate() error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.UserApikey{})
	if err != nil {
		return err
	}

	userCount, err := r.Count()
	if err != nil {
		return err
	}

	if userCount == 0 {
		firstAdminUser := model.User{
			Username:    "administrator",
			FirstName:   "BeeTimeClock",
			LastName:    "Administrator",
			AccessLevel: model.USER_ACCESS_LEVEL_ADMIN,
		}

		firstAdminUser.SetPassword("lol123")

		err = r.Insert(&firstAdminUser)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *User) FindAll() ([]model.User, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return nil, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var items []model.User

	result := db.Find(&items)
	return items, result.Error
}

func (r *User) FindByID(id uint) (model.User, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.User{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.User
	result := db.Find(&item, "id = ?", id)

	if result.RowsAffected == 0 {
		return model.User{}, fmt.Errorf("no user with id %d found", id)
	}

	return item, result.Error
}

func (r *User) FindByUsername(username string) (model.User, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.User{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.User
	result := db.Find(&item, "username = ?", username)

	if result.RowsAffected == 0 {
		return model.User{}, ErrUserNotFound
	}

	return item, result.Error
}

func (r *User) FindUserByApikey(apikey string) (model.User, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.User{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.UserApikey
	result := db.Preload(clause.Associations).Find(&item, "apikey = ?", apikey)

	if result.RowsAffected == 0 {
		return model.User{}, ErrUserNotFound
	}

	return item.User, nil
}

func (r *User) Insert(user *model.User) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(&user)
	return result.Error
}

func (r *User) Update(user *model.User) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(&user)
	return result.Error
}

func (r *User) Delete(user *model.User) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Unscoped().Delete(&user)
	return result.Error
}

func (r *User) Count() (int64, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return 0, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var count int64
	result := db.Model(&model.User{}).Count(&count)
	return count, result.Error
}

func (r *User) UserApikeyFindAll() ([]model.UserApikey, error) {
	var items []model.UserApikey
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

func (r User) UserApikeyFindById(id uint) (model.UserApikey, error) {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return model.UserApikey{}, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	var item model.UserApikey
	result := db.Find(&item, "id = ?", id)
	if result.Error != nil {
		return model.UserApikey{}, result.Error
	}

	if result.RowsAffected == 0 {
		return model.UserApikey{}, ErrUserApikeyNotFound
	}
	return item, result.Error
}

func (r User) UserApikeyInsert(push *model.UserApikey) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Create(push)
	return result.Error
}

func (r User) UserApikeyUpdate(push *model.UserApikey) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Updates(push)
	return result.Error
}

func (r User) UserApikeyDelete(push *model.UserApikey) error {
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Delete(push)
	return result.Error
}

func (r *User) UserApikeyFindAllByUserID(userID uint) ([]model.UserApikey, error) {
	var items []model.UserApikey
	db, err := r.env.DatabaseManager.GetConnection()
	if err != nil {
		return items, err
	}
	defer r.env.DatabaseManager.CloseConnection(db)

	result := db.Find(&items, "user_id = ?", userID)
	if result.Error != nil {
		return items, result.Error
	}
	return items, result.Error
}
