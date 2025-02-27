package handler

import (
	"net/http"
	"strconv"

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
}

func NewExternalWork(env *core.Environment, user *repository.User, externalWork *repository.ExternalWork) *ExternalWork {
	return &ExternalWork{
		env:          env,
		user:         user,
		externalWork: externalWork,
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

	idParam := c.Param("id")
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

	if len(externalWorkItem.WorkExpanses) == 0 {
		fromDay := helper.GetDayDate(externalWorkItem.From)
		tillDay := helper.GetDayDate(externalWorkItem.Till)

		currentDay := fromDay
		for currentDay.Before(tillDay.AddDate(0, 0, 1)) {
			expanseItem := model.ExternalWorkExpense{
				ExternalWork: externalWorkItem,
				Date:         currentDay,
			}

			externalWorkItem.WorkExpanses = append(externalWorkItem.WorkExpanses, expanseItem)

			currentDay = currentDay.AddDate(0, 0, 1)
		}
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(externalWorkItem))
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
		User:        user,
		From:        externalWorkCreateRequest.From,
		Till:        externalWorkCreateRequest.Till,
		Description: externalWorkCreateRequest.Description,
		Identifier:  uuid.New(),
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
