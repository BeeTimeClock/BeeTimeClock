package handler

import (
	"net/http"
	"strconv"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type Holiday struct {
	env     *core.Environment
	holiday *repository.Holiday
}

func NewHoliday(env *core.Environment, holiday *repository.Holiday) *Holiday {
	return &Holiday{
		env:     env,
		holiday: holiday,
	}
}

func (h *Holiday) HolidayYearGet(c *gin.Context) {
	year, err := strconv.ParseInt(c.Param("year"), 10, 64)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	holidays, err := h.holiday.HolidayFindByYear(int(year))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(holidays))
}
