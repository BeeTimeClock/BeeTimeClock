package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/microsoft"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

const sessionVarUser = "user"
const sessionVarIsAdministrator = "is_administrator"

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

	switch {
	case strings.HasPrefix(authorizationHeader, "Bearer "):
		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
		c.Set("token", tokenString)

		switch strings.ToLower(authProvider) {
		case "local", "":
			a.localAuthRequired(c, tokenString)
		case "microsoft":
			a.microsoftAuthRequired(c, tokenString)
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(fmt.Errorf("auth provider %s not supported", authProvider)))
			return
		}
		return
	case strings.HasPrefix(authorizationHeader, "Apikey "):
		apikey := strings.Replace(authorizationHeader, "Apikey ", "", 1)

		user, err := a.user.FindUserByApikey(apikey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(fmt.Errorf("no access rights")))
			return
		}

		c.Set(sessionVarUser, user)
		c.Set(sessionVarIsAdministrator, user.AccessLevel == model.USER_ACCESS_LEVEL_ADMIN)
		c.Next()
		return
	default:
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(fmt.Errorf("no supported header")))
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

func IsAdministrator(c *gin.Context) bool {
	return c.GetBool(sessionVarIsAdministrator)
}

func (a *AuthProvider) AuthProviders(c *gin.Context) {
	type AuthProviders struct {
		Local     bool
		Microsoft bool
	}

	hasMicrosoft := microsoft.GetMicrosoftClientId() != ""

	authProviders := AuthProviders{
		Local:     true,
		Microsoft: hasMicrosoft,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(authProviders))
}
