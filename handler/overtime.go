package handler

import (
	"net/http"
	"strconv"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/BeeTimeClock/BeeTimeClock-Server/worker"
	"github.com/gin-gonic/gin"
)

type Overtime struct {
	env            *core.Environment
	overtime       *repository.Overtime
	user           *repository.User
	overtimeWorker *worker.Overtime
}

func NewOvertime(env *core.Environment, user *repository.User, overtime *repository.Overtime, overtimeWorker *worker.Overtime) *Overtime {
	return &Overtime{
		env:            env,
		user:           user,
		overtime:       overtime,
		overtimeWorker: overtimeWorker,
	}
}

func (h *Overtime) OvertimeGetAll(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	overtimeMonths, err := h.overtime.OvertimeMonthQuotaFindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(overtimeMonths))
}

func (h *Overtime) OvertimeTotal(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	overtimeMonths, err := h.overtime.OvertimeMonthQuotaFindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	result := model.SumResult{
		Total: 0,
	}
	for _, overtimeQuota := range overtimeMonths {
		result.Total += overtimeQuota.Hours
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Overtime) OvertimeCurrentUserCalculateMonth(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	monthParam := c.Param("month")
	month, err := strconv.Atoi(monthParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	result, created, err := h.overtimeWorker.CalculateMonth(user.ID, year, month)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if created {
		c.JSON(http.StatusCreated, model.NewSuccessResponse(result))
		return
	}
	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}
