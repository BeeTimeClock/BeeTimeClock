package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/BeeTimeClock/BeeTimeClock-Server/microsoft"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
)

func (a *AuthProvider) microsoftAuthRequired(c *gin.Context, tokenString string) {
	token, err := verifyToken(tokenString)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if !token.Valid {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token not valid"))
		return
	}

	username := token.Claims.(jwt.MapClaims)["preferred_username"].(string)

	user, err := a.user.FindByUsername(username)
	if err == repository.ErrUserNotFound {
		fullName := token.Claims.(jwt.MapClaims)["name"].(string)

		names := strings.Split(fullName, " ")

		user = model.NewUser(username)
		user.FirstName = strings.Join(names[:len(names)-1], " ")
		user.LastName = names[len(names)-1]

		err = a.user.Insert(&user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Set(sessionVarUser, user)
	c.Set(sessionVarIsAdministrator, user.AccessLevel == model.USER_ACCESS_LEVEL_ADMIN)
	c.Next()
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	keySet, err := jwk.Fetch(context.Background(), "https://login.microsoftonline.com/common/discovery/v2.0/keys")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwa.RS256.String() {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		keys, ok := keySet.LookupKeyID(kid)
		if !ok {
			return nil, fmt.Errorf("key %v not found", kid)
		}
		issuer, err := token.Claims.GetIssuer()
		if err != nil {
			return nil, err
		}
		issuerTemplated := strings.ReplaceAll(keys.PrivateParams()["issuer"].(string), "{tenantid}", microsoft.GetMicrosoftTenantId())
		if issuer != issuerTemplated {
			return nil, fmt.Errorf("wrong issuer")
		}

		aud, err := token.Claims.GetAudience()
		if err != nil {
			return nil, err
		}
		jwt.MarshalSingleStringAsArray = false
		audBytes, err := aud.MarshalJSON()
		if err != nil {
			return nil, err
		}

		cleanedAud := strings.Trim(strings.TrimSpace(string(audBytes)), "\"")

		if cleanedAud != microsoft.GetMicrosoftClientId() {
			return nil, fmt.Errorf("wrong client")
		}

		publickey := &rsa.PublicKey{}
		err = keys.Raw(publickey)
		if err != nil {
			return nil, fmt.Errorf("could not parse pubkey")
		}

		return publickey, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (a *AuthProvider) MicrosoftAuthSettings(c *gin.Context) {
	c.JSON(http.StatusOK, model.NewSuccessResponse(gin.H{
		"ClientID": microsoft.GetMicrosoftClientId(),
		"TenantID": microsoft.GetMicrosoftTenantId(),
	}))
}
