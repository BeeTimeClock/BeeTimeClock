package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type Timestamp struct {
	env       *core.Environment
	user      *repository.User
	timestamp *repository.Timestamp
}

func NewTimestamp(env *core.Environment, user *repository.User, timestamp *repository.Timestamp) *Timestamp {
	return &Timestamp{
		env:       env,
		user:      user,
		timestamp: timestamp,
	}
}

func (h *Timestamp) TimestampGetAll(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	timestamps, err := h.timestamp.FindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(timestamps))
}

func (h *Timestamp) TimestampQueryLast(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	lastTimestamp, err := h.timestamp.FindLastByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, lastTimestamp)
}

func (h *Timestamp) TimestampActionCheckIn(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	var timestampActionCheckInRequest model.TimestampActionCheckInRequest
	err = c.BindJSON(&timestampActionCheckInRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	timestampCount, err := h.timestamp.CountByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if timestampCount > 0 {
		lastTimestamp, err := h.timestamp.FindLastByUserID(user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}

		if !lastTimestamp.IsComplete() {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("there is an open timestamp")))
			return
		}
	}

	timestamp := model.Timestamp{
		User:            &user,
		ComingTimestamp: time.Now(),
		IsHomeoffice:    timestampActionCheckInRequest.IsHomeoffice,
	}

	err = h.timestamp.Insert(&timestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(timestamp))
}

func (h *Timestamp) TimestampActionCheckOut(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	lastTimestamp, err := h.timestamp.FindLastByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if lastTimestamp.IsComplete() {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Errorf("there is no open timestamp"))
		return
	}

	lastTimestamp.GoingTimestamp = time.Now()

	err = h.timestamp.Update(&lastTimestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(lastTimestamp))
}

func (h *Timestamp) TimestampCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	var timestampCreateRequest model.TimestampCreateRequest
	err = c.BindJSON(&timestampCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	timestamp := model.Timestamp{
		User:            &user,
		ComingTimestamp: timestampCreateRequest.ComingTimestamp,
		GoingTimestamp:  timestampCreateRequest.GoingTimestamp,
		IsHomeoffice:    timestampCreateRequest.IsHomeoffice,
	}

	err = h.timestamp.Insert(&timestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(timestamp))
}

func (h *Timestamp) TimestampQueryCurrentMonthGrouped(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	currentYear, currentMonth, _ := time.Now().Date()
	currentLocation := time.Now().Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	timestamps, err := h.timestamp.FindByUserIDAndDate(user.ID, firstOfMonth, lastOfMonth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
	}

	type TimestampGroup struct {
		Date            time.Time
		IsHomeoffice    bool
		Timestamps      []model.Timestamp
		WorkingHours    float64
		SubtractedHours float64
	}
	grouped := make(map[time.Time]TimestampGroup)

	for _, timestamp := range timestamps {
		year, month, day := timestamp.ComingTimestamp.Date()
		timestamp_date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		if _, exists := grouped[timestamp_date]; !exists {
			grouped[timestamp_date] = TimestampGroup{
				IsHomeoffice: true,
				Date:         timestamp_date,
			}
		}
		group := grouped[timestamp_date]

		if !timestamp.IsHomeoffice {
			group.IsHomeoffice = false
		}
		group.Timestamps = append(grouped[timestamp_date].Timestamps, timestamp)

		workingHours, subtractedHours := timestamp.CalculateWorkingHours()
		group.WorkingHours += workingHours
		group.SubtractedHours += subtractedHours

		grouped[timestamp_date] = group
	}

	result := []TimestampGroup{}
	for _, value := range grouped {
		result = append(result, value)
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}
