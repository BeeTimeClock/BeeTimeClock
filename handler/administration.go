package handler

import (
	"net/http"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type Administration struct {
	env      *core.Environment
	settings *repository.Settings
}

func NewAdministration(env *core.Environment, settings *repository.Settings) *Administration {
	return &Administration{
		env:      env,
		settings: settings,
	}
}

func (h Administration) AdministrationGetSettings(c *gin.Context) {
	settings, err := h.settings.SettingsFind()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(settings))
}

func (h Administration) AdministrationUpdateSettings(c *gin.Context) {
	var settings model.Settings

	err := c.BindJSON(&settings)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	h.settings.SettingsUpdate(&settings)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(settings))
}
