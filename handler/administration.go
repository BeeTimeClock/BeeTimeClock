package handler

import (
	"net/http"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/BeeTimeClock/BeeTimeClock-Server/worker"
	"github.com/gin-gonic/gin"
)

type Administration struct {
	env      *core.Environment
	settings *repository.Settings
	absence  *repository.Absence
}

func NewAdministration(env *core.Environment, settings *repository.Settings, absence *repository.Absence) *Administration {
	return &Administration{
		env:      env,
		settings: settings,
		absence:  absence,
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

func (h Administration) AdministrationNotifyAbsenceWeek(c *gin.Context) {
	worker.NotifyAbsenceWeek(h.env, h.absence)
	c.Status(http.StatusNoContent)
}
