package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/BeeTimeClock/BeeTimeClock-Server/worker"
	"github.com/gin-gonic/gin"
)

type Timestamp struct {
	env             *core.Environment
	user            *repository.User
	timestamp       *repository.Timestamp
	absence         *repository.Absence
	settings        *repository.Settings
	holiday         *repository.Holiday
	timestampWorker *worker.Timestamp
}

func NewTimestamp(env *core.Environment, user *repository.User, timestamp *repository.Timestamp, absence *repository.Absence, settings *repository.Settings, holiday *repository.Holiday, timestampWorker *worker.Timestamp) *Timestamp {
	return &Timestamp{
		env:             env,
		user:            user,
		timestamp:       timestamp,
		absence:         absence,
		settings:        settings,
		holiday:         holiday,
		timestampWorker: timestampWorker,
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

func (h *Timestamp) TimestampQuerySuspicious(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	timestamps, err := h.timestamp.FindSuspiciousTimestampsByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(timestamps))
}

func (h *Timestamp) TimestampQuerySuspiciousCount(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	timestamps, err := h.timestamp.FindSuspiciousTimestampsByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(model.CountResult{
		Count: len(timestamps),
	}))
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

	isHomeoffice, err := h.isHomeoffice(c, timestampActionCheckInRequest.IsHomeoffice)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	timestamp := model.Timestamp{
		User:            &user,
		ComingTimestamp: time.Now(),
		IsHomeoffice:    isHomeoffice,
	}

	err = h.timestamp.Insert(&timestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(timestamp))
}

func (h *Timestamp) isHomeoffice(c *gin.Context, prefered bool) (bool, error) {
	settings, err := h.settings.SettingsFind()
	if err != nil {
		return false, err
	}

	if *settings.CheckinDetectionByIPAddress == true {
		clientIps, err := getClientIPByHeaders(c)
		if err != nil {
			return false, err
		}

		isOfficeIp := false

		for _, input := range strings.Split(clientIps, ",") {
			clientIp := strings.TrimSpace(input)
			fmt.Printf("Use ClientIP Detection: %s\n", clientIps)
			for _, officeIp := range settings.OfficeIPAddresses {
				if officeIp.IPAddress == clientIp {
					isOfficeIp = true
					break
				}
			}

			if isOfficeIp {
				break
			}
		}

		return !isOfficeIp, nil
	} else {
		return prefered, nil
	}
}

func (h *Timestamp) TimestampActionCheckOut(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	var timestampCheckoutActionRequest model.TimestampActionCheckoutRequest
	err = c.BindJSON(&timestampCheckoutActionRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	lastTimestamp, err := h.timestamp.FindLastByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if lastTimestamp.IsComplete() {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("there is no open timestamp")))
		return
	}

	isHomeoffice, err := h.isHomeoffice(c, timestampCheckoutActionRequest.IsHomeoffice)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	lastTimestamp.GoingTimestamp = time.Now()
	lastTimestamp.IsHomeofficeGoing = isHomeoffice

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

	timestampCorrection := model.TimestampCorrection{
		Timestamp:    timestamp,
		ChangeReason: timestampCreateRequest.ChangeReason,
	}

	err = h.timestamp.TimestampCorrectionInsert(&timestampCorrection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(timestamp))
}

func (h *Timestamp) TimestampUserQueryMonthGrouped(c *gin.Context) {
	user, success := getUserFromParam(c, h.user)
	if !success {
		return
	}

	year, month, success := getYearMonthFromParam(c)
	if !success {
		return
	}

	result, err := h.timestampWorker.CalculateMonth(user.ID, year, month)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Timestamp) TimestampUserQueryMonthOvertime(c *gin.Context) {
	user, success := getUserFromParam(c, h.user)
	if !success {
		return
	}

	year, month, success := getYearMonthFromParam(c)
	if !success {
		return
	}

	overtimeHours, err := h.overtimeMonth(user.ID, year, month)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	holidays, err := h.holiday.HolidayFindByYear(year)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	neededHours := model.GetNeededHoursForMonth(holidays, year, month)

	subtractedHours := 0.0
	if overtimeHours > 0 {
		switch user.OvertimeSubtractionModel {
		case model.OVERTIME_SUBTRACTION_MODEL_HOURS:
			subtractedHours = user.OvertimeSubtractionAmount

			if subtractedHours > overtimeHours {
				subtractedHours = overtimeHours
			}

			break
		case model.OVERTIME_SUBTRACTION_MODEL_PERCENTAGE:
			subtractedHours = neededHours / 100 * user.OvertimeSubtractionAmount
			if subtractedHours > overtimeHours {
				subtractedHours = overtimeHours
			}
			break
		}
	}

	result := model.OvertimeResult{
		Total:      overtimeHours,
		Needed:     neededHours,
		Subtracted: subtractedHours,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Timestamp) TimestampCurrentUserQueryCurrentMonthGrouped(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	currentYear, currentMonth, _ := time.Now().Date()
	result, err := h.timestampWorker.CalculateMonth(user.ID, currentYear, int(currentMonth))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Timestamp) TimestampCurrentUserQueryCurrentMonthOvertime(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	overtimeHours, err := h.overtimeCurrentMonth(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	result := model.SumResult{
		Total: overtimeHours,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Timestamp) TimestampQueryMonthGrouped(c *gin.Context) {
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

	result, err := h.timestampWorker.CalculateMonth(user.ID, year, month)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result.TimestampGroups))
}

func (h *Timestamp) TimestampQueryMonthOvertime(c *gin.Context) {
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

	overtimeHours, err := h.overtimeMonth(user.ID, year, month)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	holidays, err := h.holiday.HolidayFindByYear(year)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	neededHours := model.GetNeededHoursForMonth(holidays, year, month)

	subtractedHours := 0.0
	if overtimeHours > 0 {
		switch user.OvertimeSubtractionModel {
		case model.OVERTIME_SUBTRACTION_MODEL_HOURS:
			subtractedHours = user.OvertimeSubtractionAmount

			if subtractedHours > overtimeHours {
				subtractedHours = overtimeHours
			}

			break
		case model.OVERTIME_SUBTRACTION_MODEL_PERCENTAGE:
			subtractedHours = neededHours / 100 * user.OvertimeSubtractionAmount
			if subtractedHours > overtimeHours {
				subtractedHours = overtimeHours
			}
			break
		}
	}

	result := model.OvertimeResult{
		Total:      overtimeHours,
		Needed:     neededHours,
		Subtracted: subtractedHours,
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Timestamp) overtimeCurrentMonth(userID uint) (float64, error) {
	currentYear, currentMonth, _ := time.Now().Date()
	return h.overtimeMonth(userID, currentYear, int(currentMonth))
}

func (h *Timestamp) overtimeMonth(userID uint, year int, month int) (float64, error) {
	result, err := h.timestampWorker.CalculateMonth(userID, year, month)
	if err != nil {
		return 0.0, err
	}

	return result.OvertimeHours, nil
}

func (h *Timestamp) TimestampCorrectionCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	timestampIdParam := c.Param("timestampID")
	timestampId, err := strconv.ParseUint(timestampIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	var timestampCorrectionCreateRequest model.TimestampCorrectionCreateRequest
	err = c.BindJSON(&timestampCorrectionCreateRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	timestamp, err := h.timestamp.FindByID(uint(timestampId))
	if err != nil {
		if err == repository.ErrTimestampNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, model.NewErrorResponse(err))
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		}

		return
	}

	if timestamp.UserID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	timestampCorrection := model.TimestampCorrection{
		Timestamp:          timestamp,
		ChangeReason:       timestampCorrectionCreateRequest.ChangeReason,
		OldComingTimestamp: timestamp.ComingTimestamp,
		OldGoingTimestamp:  timestamp.GoingTimestamp,
	}

	err = h.timestamp.TimestampCorrectionInsert(&timestampCorrection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	timestamp.ComingTimestamp = timestampCorrectionCreateRequest.NewComingTimestamp
	timestamp.GoingTimestamp = timestampCorrectionCreateRequest.NewGoingTimestamp
	timestamp.IsHomeoffice = timestampCorrectionCreateRequest.IsHomeoffice

	err = h.timestamp.Update(&timestamp)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(timestamp))
}

func (h *Timestamp) queryMonths(c *gin.Context, userID uint) (map[int][]int, bool) {
	yearMonths, err := h.timestamp.FindYearMonthsWithTimestampsByUserId(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return nil, false
	}

	result := make(map[int][]int)

	for _, group := range yearMonths {
		if _, exists := result[group.Year]; !exists {
			result[group.Year] = []int{}
		}

		result[group.Year] = append(result[group.Year], group.Month)
	}

	currentYear := time.Now().Year()

	if _, exists := result[currentYear]; !exists {
		result[currentYear] = []int{}
	}

	currentMonth := int(time.Now().Month())

	contains := false
	for _, month := range result[currentYear] {
		if month == currentMonth {
			contains = true
			break
		}
	}
	if !contains {
		result[currentYear] = append(result[currentYear], currentMonth)
	}

	return result, true
}

func (h *Timestamp) TimestampQueryMonths(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	result, success := h.queryMonths(c, user.ID)
	if !success {
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Timestamp) TimestampUserQueryMonths(c *gin.Context) {
	user, success := getUserFromParam(c, h.user)
	if !success {
		return
	}

	result, success := h.queryMonths(c, user.ID)
	if !success {
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}
