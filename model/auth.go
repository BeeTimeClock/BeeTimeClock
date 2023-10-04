package model

import "github.com/golang-jwt/jwt/v4"

type AuthRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type AuthResponse struct {
	Token string
}

type AuthInfo struct {
	jwt.RegisteredClaims
	UserID uint
}
