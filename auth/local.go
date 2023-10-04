package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (a *AuthProvider) localAuthRequired(c *gin.Context, tokenString string) {
	var authInfo model.AuthInfo
	tkn, err := jwt.ParseWithClaims(tokenString, &authInfo, func(token *jwt.Token) (interface{}, error) {
		return a.env.Secret, nil
	})

	if err != nil {
		fmt.Println(err)
		if err == jwt.ErrSignatureInvalid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(fmt.Errorf("Signature not valid")))
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	if !tkn.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	user, err := a.user.FindByID(authInfo.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Set(sessionVarUser, user)
	c.Next()
}

func (a *AuthProvider) Auth(c *gin.Context) {
	var authRequest model.AuthRequest

	err := c.BindQuery(&authRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	user, err := a.user.FindByUsername(authRequest.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	valid, err := user.CheckPassword(authRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	if !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(fmt.Errorf("username or password wrong")))
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)

	authSession := model.AuthInfo{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authSession)
	tokenString, err := token.SignedString(a.env.Secret)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(model.AuthResponse{
		Token: tokenString,
	}))
}
