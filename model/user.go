package model

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserWorkTimeModelNotFoundForTimestamp = errors.New("no worktime model found for this timestamp")

// MaxRingtoneTones is the maximum number of tones a coming/going ringtone (RTTTL) may contain.
const MaxRingtoneTones = 24

var ErrRingtoneTooManyTones = fmt.Errorf("ringtone must not contain more than %d tones", MaxRingtoneTones)

// RingtoneToneCount returns the number of tones in an RTTTL string. The tones
// are the comma-separated entries of the last colon-separated section.
func RingtoneToneCount(rtttl string) int {
	rtttl = strings.TrimSpace(rtttl)
	if rtttl == "" {
		return 0
	}

	parts := strings.Split(rtttl, ":")
	notes := parts[len(parts)-1]

	count := 0
	for _, note := range strings.Split(notes, ",") {
		if strings.TrimSpace(note) != "" {
			count++
		}
	}
	return count
}

// Validate checks the update request for invalid values.
func (r *UserUpdateRequest) Validate() error {
	if RingtoneToneCount(r.ComingRingtone) > MaxRingtoneTones {
		return ErrRingtoneTooManyTones
	}
	if RingtoneToneCount(r.GoingRingtone) > MaxRingtoneTones {
		return ErrRingtoneTooManyTones
	}
	return nil
}

type UserAccessLevel string
type UserTokenType string

const (
	USER_ACCESS_LEVEL_ADMIN UserAccessLevel = "admin"
	USER_ACCESS_LEVEL_USER  UserAccessLevel = "user"
	USER_TOKEN_TYPE_CHIP    UserTokenType   = "chip"
)

type User struct {
	gorm.Model
	Username            string `gorm:"unique"`
	Password            string `json:"-"`
	FirstName           string
	LastName            string
	AccessLevel         UserAccessLevel
	HolidayDaysPerYear  uint
	WorkingHoursPerWeek float64
	StaffNumber         int64
	WorkTimeModels      []UserWorktime
	AllowGravatar       bool
	ComingRingtone      string
	GoingRingtone       string
}

func NewUser(username string) User {
	return User{
		Username:            username,
		WorkingHoursPerWeek: 38.0,
		AccessLevel:         USER_ACCESS_LEVEL_USER,
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
	StaffNumber         int64
}

type UserUpdateRequest struct {
	AccessLevel         UserAccessLevel
	FirstName           string
	LastName            string
	HolidayDaysPerYear  uint
	WorkingHoursPerWeek float64
	StaffNumber         int64
	AllowGravatar       bool
	ComingRingtone      string
	GoingRingtone       string
}

type UserResponse struct {
	gorm.Model
	Username       string
	FirstName      string
	LastName       string
	AccessLevel    string
	StaffNumber    int64
	AllowGravatar  bool
	ComingRingtone string
	GoingRingtone  string
}

type UserApikey struct {
	gorm.Model
	UserID      uint
	Description string
	User        User
	Apikey      string `gorm:"unique"`
	ValidTill   time.Time
}

type UserToken struct {
	gorm.Model
	UserID          uint
	User            User
	TokenType       UserTokenType
	TokenIdentifier string `gorm:"unique"`
}

type UserApikeyCreateRequest struct {
	Description string `binding:"required"`
	ValidTill   time.Time
}

type UserTokenCreateRequest struct {
	TokenType       UserTokenType
	TokenIdentifier string `binding:"required"`
}

type UserApikeyResponse struct {
	gorm.Model
	Description string
	ValidTill   time.Time
}

func (u *User) GetUserResponse() UserResponse {
	return UserResponse{
		Model:         u.Model,
		Username:      u.Username,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		AccessLevel:    string(u.AccessLevel),
		StaffNumber:    u.StaffNumber,
		AllowGravatar:  u.AllowGravatar,
		ComingRingtone: u.ComingRingtone,
		GoingRingtone:  u.GoingRingtone,
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

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) GetWorkTimeModel(date time.Time) (*WorkTimeModel, error) {
	idx := slices.IndexFunc(u.WorkTimeModels, func(i UserWorktime) bool {
		return i.ValidFrom.Before(date) && (i.ValidTill == nil || i.ValidTill.After(date))
	})

	if idx == -1 {
		return nil, ErrUserWorkTimeModelNotFoundForTimestamp
	}

	return &u.WorkTimeModels[idx].WorkTimeModel, nil
}

type UserWorktime struct {
	gorm.Model
	UserID          uint
	User            User
	WorkTimeModelID uint
	WorkTimeModel   WorkTimeModel
	ValidFrom       time.Time
	ValidTill       *time.Time
}

type UserWorktimeCreateRequest struct {
	WorktimeModelID uint      `binding:"required"`
	ValidFrom       time.Time `binding:"required"`
	ValidTill       *time.Time
}

type UserWorktimeUpdateRequest struct {
	ValidFrom time.Time
	ValidTill *time.Time
}
