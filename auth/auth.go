package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

const sessionVarUser = "user"

type AuthHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

type AuthProvider struct {
	env  *core.Environment
	user *repository.User
}

func NewAuthProvider(env *core.Environment, user *repository.User) AuthProvider {
	return AuthProvider{
		env:  env,
		user: user,
	}
}

func (a *AuthProvider) AuthRequired(c *gin.Context) {
	authProvider := c.GetHeader("X-Auth-Provider")
	authorizationHeader := c.GetHeader("Authorization")
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)

	switch strings.ToLower(authProvider) {
	case "local", "":
		a.localAuthRequired(c, tokenString)
	case "microsoft":
		a.microsoftAuthRequired(c, tokenString)
	default:
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(fmt.Errorf("auth provider %s not supported", authProvider)))
		return
	}
}

func AdministratorAccessRequired(c *gin.Context) {
	user, err := GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	if user.AccessLevel != model.USER_ACCESS_LEVEL_ADMIN {
		c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(fmt.Errorf("no access rights")))
		return
	}

	c.Next()
}

func GetUserFromSession(c *gin.Context) (model.User, error) {
	user, exists := c.Get(sessionVarUser)
	if !exists {
		return model.User{}, fmt.Errorf("no user in session")
	}

	return user.(model.User), nil
}

func (a *AuthProvider) AuthProviders(c *gin.Context) {
	type AuthProviders struct {
		Local     bool
		Microsoft bool
	}

	hasMicrosoft := getMicrosoftClientId() != ""

	authProviders := AuthProviders{
		Local:     true,
		Microsoft: hasMicrosoft,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(authProviders))
}
