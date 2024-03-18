package handler

import (
	"net/http"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type Migration struct {
	env       *core.Environment
	migration *repository.Migration
}

func NewMigration(env *core.Environment, migration *repository.Migration) *Migration {
	return &Migration{
		env:       env,
		migration: migration,
	}
}

func (h *Migration) AdministrationMigrationGetAll(c *gin.Context) {
	migrations, err := h.migration.MigrationFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(migrations))
}
