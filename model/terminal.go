package model

import (
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Terminal struct {
	gorm.Model
	TerminalName string `gorm:"uniqueIndex"`
	ClientId     string `gorm:"uniqueIndex"`
	Apikey       string `gorm:"not null"`
}

func (t *Terminal) CheckApikey(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(t.Apikey), []byte(plaintext))
	return err == nil, err
}

func (t *Terminal) GenerateApikey() (string, error) {
	apikeyPlain := helper.RandomString(32)

	bytes, err := bcrypt.GenerateFromPassword([]byte(apikeyPlain), 14)
	if err != nil {
		return "", err
	}

	t.Apikey = string(bytes)
	return apikeyPlain, nil
}

type TerminalCreateRequest struct {
	TerminalName string `binding:"required"`
}

type TerminalCheckinRequest struct {
	TokenIdentifier string `binding:"required"`
}

type TerminalCheckoutRequest struct {
	TokenIdentifier string `binding:"required"`
}
