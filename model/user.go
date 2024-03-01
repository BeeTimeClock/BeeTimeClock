package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserAccessLevel string
type OvertimeSubtractionModel string

const (
	USER_ACCESS_LEVEL_ADMIN               UserAccessLevel          = "admin"
	USER_ACCESS_LEVEL_USER                UserAccessLevel          = "user"
	OVERTIME_SUBTRACTION_MODEL_HOURS      OvertimeSubtractionModel = "hours"
	OVERTIME_SUBTRACTION_MODEL_PERCENTAGE OvertimeSubtractionModel = "percentage"
)

type User struct {
	gorm.Model
	Username                  string `gorm:"unique"`
	Password                  string
	FirstName                 string
	LastName                  string
	AccessLevel               UserAccessLevel
	HolidayDaysPerYear        uint
	WorkingHoursPerWeek       float64
	OvertimeSubtractionModel  OvertimeSubtractionModel
	OvertimeSubtractionAmount float64
}

func NewUser(username string) User {
	return User{
		Username:                  username,
		HolidayDaysPerYear:        30,
		WorkingHoursPerWeek:       38.0,
		AccessLevel:               USER_ACCESS_LEVEL_USER,
		OvertimeSubtractionModel:  OVERTIME_SUBTRACTION_MODEL_PERCENTAGE,
		OvertimeSubtractionAmount: 10,
	}
}

type UserDeleteQuery struct {
	UserID uint `binding:"required"`
}

type UserCreateRequest struct {
	Username            string          `binding:"required"`
	Password            string          `binding:"required"`
	AccessLevel         UserAccessLevel `binding:"required"`
	FirstName           string
	LastName            string
	HolidayDaysPerYear  uint
	WorkingHoursPerWeek float64
}

type UserUpdateRequest struct {
	AccessLevel               UserAccessLevel
	FirstName                 string
	LastName                  string
	HolidayDaysPerYear        uint
	WorkingHoursPerWeek       float64
	OvertimeSubtractionAmount float64
	OvertimeSubtractionModel  OvertimeSubtractionModel
}

type UserResponse struct {
	gorm.Model
	Username                  string
	FirstName                 string
	LastName                  string
	AccessLevel               string
	OvertimeSubtractionModel  OvertimeSubtractionModel
	OvertimeSubtractionAmount float64
}

type UserApikey struct {
	gorm.Model
	UserID      uint
	Description string
	User        User
	Apikey      string `gorm:"unique"`
	ValidTill   time.Time
}

type UserApikeyCreateRequest struct {
	Description string `binding:"required"`
	ValidTill   time.Time
}

type UserApikeyResponse struct {
	gorm.Model
	Description string
	ValidTill   time.Time
}

func (u *User) GetUserResponse() UserResponse {
	return UserResponse{
		Model:                     u.Model,
		Username:                  u.Username,
		FirstName:                 u.FirstName,
		LastName:                  u.LastName,
		AccessLevel:               string(u.AccessLevel),
		OvertimeSubtractionModel:  u.OvertimeSubtractionModel,
		OvertimeSubtractionAmount: u.OvertimeSubtractionAmount,
	}
}

func (ua *UserApikey) GetUserApikeyResponse() UserApikeyResponse {
	return UserApikeyResponse{
		Model:       ua.Model,
		Description: ua.Description,
		ValidTill:   ua.ValidTill,
	}
}

func (u *User) CheckPassword(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plaintext))
	return err == nil, err
}

func (u *User) SetPassword(plaintext string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), 14)
	if err != nil {
		return err
	}

	u.Password = string(bytes)
	return nil
}
