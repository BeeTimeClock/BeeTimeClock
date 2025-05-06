package handler

import (
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"github.com/BeeTimeClock/BeeTimeClock-Server/microsoft"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExternalWork struct {
	env          *core.Environment
	user         *repository.User
	externalWork *repository.ExternalWork
	holiday      *repository.Holiday
}

func NewExternalWork(env *core.Environment, user *repository.User, externalWork *repository.ExternalWork, holiday *repository.Holiday) *ExternalWork {
	return &ExternalWork{
		env:          env,
		user:         user,
		externalWork: externalWork,
		holiday:      holiday,
	}
}

func (h *ExternalWork) ExternalWorkGetAll(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	externalWorkItems, err := h.externalWork.ExternalWorkFindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkItems))
}

func (h *ExternalWork) ExternalWorkGetById(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	idParam := c.Param("externalWorkId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkItem, err := h.externalWork.ExternalWorkFindById(uint(id), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if externalWorkItem.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	fromDay := helper.GetDayDate(externalWorkItem.From)
	tillDay := helper.GetDayDate(externalWorkItem.Till)

	currentDay := fromDay
	for currentDay.Before(tillDay.AddDate(0, 0, 1)) {
		if !slices.ContainsFunc(externalWorkItem.WorkExpanses, func(n model.ExternalWorkExpense) bool {
			return n.Date.Round(24*time.Hour).UTC() == currentDay
		}) {
			expanseItem := model.ExternalWorkExpense{
				ExternalWork:   externalWorkItem,
				ExternalWorkID: externalWorkItem.ID,
				Date:           currentDay,
			}

			externalWorkItem.WorkExpanses = append(externalWorkItem.WorkExpanses, expanseItem)
		}

		currentDay = currentDay.AddDate(0, 0, 1)
	}

	holidays, err := h.holiday.HolidayFindByDateRange(fromDay, tillDay)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkItem.Calculate(holidays)))
}

func (h *ExternalWork) ExternalWorkCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}
	var externalWorkCreateRequest model.ExternalWorkCreateRequest

	err = c.BindJSON(&externalWorkCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWork := model.ExternalWork{
		User:                       user,
		ExternalWorkCompensationID: externalWorkCreateRequest.ExternalWorkCompensationID,
		From:                       externalWorkCreateRequest.From,
		Till:                       externalWorkCreateRequest.Till,
		Description:                externalWorkCreateRequest.Description,
		Identifier:                 uuid.New(),
	}

	err = h.externalWork.ExternalWorkInsert(&externalWork)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if microsoft.IsMicrosoftConnected() {
		eventId, err := microsoft.CreateCalendarEntryFromExternalWork(user.Username, &externalWork)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}

		externalWorkExternalEvent := model.ExternalWorkExternalEvent{
			ExternalWork:          externalWork,
			ExternalEventID:       eventId,
			ExternalEventProvider: model.EXTERNAL_EVENT_PROVIDER_MICROSOFT,
		}

		err = h.externalWork.ExternalWorkExternalEventInsert(&externalWorkExternalEvent)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(externalWork))
}

func (h *ExternalWork) ExternalWorkExpanseCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	idParam := c.Param("externalWorkId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	var externalWorkExpenseCreateRequest model.ExternalWorkExpenseCreateRequest
	err = c.BindJSON(&externalWorkExpenseCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkItem, err := h.externalWork.ExternalWorkFindById(uint(id), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if externalWorkItem.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	externalWorkExpense := model.ExternalWorkExpense{
		ExternalWork:           externalWorkItem,
		Date:                   externalWorkExpenseCreateRequest.Date,
		DepartureTime:          externalWorkExpenseCreateRequest.DepartureTime,
		ArrivalTime:            externalWorkExpenseCreateRequest.ArrivalTime,
		TravelDurationHours:    externalWorkExpenseCreateRequest.TravelDurationHours,
		PauseDurationHours:     externalWorkExpenseCreateRequest.PauseDurationHours,
		OnSiteFrom:             externalWorkExpenseCreateRequest.OnSiteFrom,
		OnSiteTill:             externalWorkExpenseCreateRequest.OnSiteTill,
		Place:                  externalWorkExpenseCreateRequest.Place,
		TravelWithPrivateCarKm: externalWorkExpenseCreateRequest.TravelWithPrivateCarKm,
	}

	err = h.externalWork.ExternalWorkExpenseInsert(&externalWorkExpense)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, externalWorkExpense)
}

func (h *ExternalWork) ExternalWorkExpanseUpdate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	idParam := c.Param("externalWorkId")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	expanseIdParam := c.Param("externalWorkExpanseId")
	expanseId, err := strconv.ParseInt(expanseIdParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	var externalWorkExpenseUpdateRequest model.ExternalWorkExpenseUpdateRequest
	err = c.BindJSON(&externalWorkExpenseUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	externalWorkItem, err := h.externalWork.ExternalWorkFindById(uint(id), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if externalWorkItem.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	externalWorkExpenseItem, err := h.externalWork.ExternalWorkExpenseFindById(uint(expanseId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	externalWorkExpenseItem.DepartureTime = externalWorkExpenseUpdateRequest.DepartureTime
	externalWorkExpenseItem.ArrivalTime = externalWorkExpenseUpdateRequest.ArrivalTime
	externalWorkExpenseItem.TravelDurationHours = externalWorkExpenseUpdateRequest.TravelDurationHours
	externalWorkExpenseItem.PauseDurationHours = externalWorkExpenseUpdateRequest.PauseDurationHours
	externalWorkExpenseItem.OnSiteFrom = externalWorkExpenseUpdateRequest.OnSiteFrom
	externalWorkExpenseItem.OnSiteTill = externalWorkExpenseUpdateRequest.OnSiteTill
	externalWorkExpenseItem.Place = externalWorkExpenseUpdateRequest.Place
	externalWorkExpenseItem.TravelWithPrivateCarKm = externalWorkExpenseUpdateRequest.TravelWithPrivateCarKm

	err = h.externalWork.ExternalWorkExpenseUpdate(&externalWorkExpenseItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, externalWorkExpenseItem)
}

func (h *ExternalWork) AdministrationExternalWorkCompensationGetAll(c *gin.Context) {
	compensations, err := h.externalWork.ExternalWorkCompensationFindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(compensations))
}

func (h *ExternalWork) AdministrationExternalWorkCompensationCreate(c *gin.Context) {
	var externalWorkCompensation model.ExternalWorkCompensation

	err := c.BindJSON(&externalWorkCompensation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	err = h.externalWork.ExternalWorkCompensationInsert(&externalWorkCompensation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkCompensation))
}
