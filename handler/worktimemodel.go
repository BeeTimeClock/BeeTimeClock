package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type WorkTimeModel struct {
	env           *core.Environment
	user          *repository.User
	workTimeModel *repository.WorkTimeModel
}

func NewWorkTimeModel(env *core.Environment, user *repository.User, workTimeModel *repository.WorkTimeModel) *WorkTimeModel {
	return &WorkTimeModel{
		env:           env,
		user:          user,
		workTimeModel: workTimeModel,
	}
}

func (h *WorkTimeModel) AdministrationWorkTimeModelGet(c *gin.Context) {
	workTimeModels, err := h.workTimeModel.WorkTimeModelFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(workTimeModels))
}

func (h *WorkTimeModel) AdministrationWorkTimeModelCreate(c *gin.Context) {
	var request model.WorkTimeModelCreateRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	workTimeModel := model.WorkTimeModel{
		Name:                      request.Name,
		WorkingHoursPerWeekday:    request.WorkingHoursPerWeekday,
		OvertimeSubtractionModel:  request.OvertimeSubtractionModel,
		OvertimeSubtractionAmount: request.OvertimeSubtractionAmount,
		HoursPerWeekdayException:  request.HoursPerWeekdayException,
		HolidayDaysPerYear:        request.HolidayDaysPerYear,
		OvertimeWarningThreshold:  request.OvertimeWarningThreshold,
	}

	err = h.workTimeModel.WorkTimeModelInsert(&workTimeModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(workTimeModel))
}

func (h *WorkTimeModel) AdministrationWorkTimeModelUpdate(c *gin.Context) {
	workTimeModelIDParam := c.Param("workTimeModelID")
	workTimeModelID, err := strconv.ParseUint(workTimeModelIDParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("invalid workTimeModelID")))
		return
	}

	workTimeModel, err := h.workTimeModel.WorkTimeModelFindById(uint(workTimeModelID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
		return
	}

	var request model.WorkTimeModelUpdateRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	workTimeModel.WorkingHoursPerWeekday = request.WorkingHoursPerWeekday
	workTimeModel.OvertimeSubtractionModel = request.OvertimeSubtractionModel
	workTimeModel.OvertimeSubtractionAmount = request.OvertimeSubtractionAmount
	workTimeModel.HoursPerWeekdayException = request.HoursPerWeekdayException
	workTimeModel.HolidayDaysPerYear = request.HolidayDaysPerYear
	workTimeModel.OvertimeWarningThreshold = request.OvertimeWarningThreshold

	err = h.workTimeModel.WorkTimeModelUpdate(&workTimeModel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(workTimeModel))
}

func (h *WorkTimeModel) AdministrationUserWorktimeGet(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("invalid userID")))
		return
	}

	userWorktimes, err := h.workTimeModel.UserWorktimeFindByUserID(uint(userID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(userWorktimes))
}

func (h *WorkTimeModel) AdministrationUserWorktimeCreate(c *gin.Context) {
	userIDParam := c.Param("userID")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("invalid userID")))
		return
	}

	var request model.UserWorktimeCreateRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	_, err = h.workTimeModel.WorkTimeModelFindById(request.WorktimeModelID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(fmt.Errorf("worktime model not found")))
		return
	}

	userWorktime := model.UserWorktime{
		UserID:          uint(userID),
		WorkTimeModelID: request.WorktimeModelID,
		ValidFrom:       request.ValidFrom,
		ValidTill:       request.ValidTill,
	}

	err = h.workTimeModel.UserWorktimeInsert(&userWorktime)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	userWorktime, err = h.workTimeModel.UserWorktimeFindByID(userWorktime.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(userWorktime))
}

func (h *WorkTimeModel) AdministrationUserWorktimeUpdate(c *gin.Context) {
	userWorktimeIDParam := c.Param("userWorktimeID")
	userWorktimeID, err := strconv.ParseUint(userWorktimeIDParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("invalid userWorktimeID")))
		return
	}

	userWorktime, err := h.workTimeModel.UserWorktimeFindByID(uint(userWorktimeID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
		return
	}

	var request model.UserWorktimeUpdateRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	if !request.ValidFrom.IsZero() {
		userWorktime.ValidFrom = request.ValidFrom
	}
	userWorktime.ValidTill = request.ValidTill

	err = h.workTimeModel.UserWorktimeUpdate(&userWorktime)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(userWorktime))
}

func (h *WorkTimeModel) AdministrationUserWorktimeDelete(c *gin.Context) {
	userWorktimeIDParam := c.Param("userWorktimeID")
	userWorktimeID, err := strconv.ParseUint(userWorktimeIDParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("invalid userWorktimeID")))
		return
	}

	userWorktime, err := h.workTimeModel.UserWorktimeFindByID(uint(userWorktimeID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
		return
	}

	err = h.workTimeModel.UserWorktimeDelete(&userWorktime)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
