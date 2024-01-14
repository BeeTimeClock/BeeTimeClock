package handler

import (
	"net/http"
	"time"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type Absence struct {
	env     *core.Environment
	user    *repository.User
	absence *repository.Absence
}

func NewAbsence(env *core.Environment, user *repository.User, absence *repository.Absence) *Absence {
	return &Absence{
		env:     env,
		user:    user,
		absence: absence,
	}
}

func (h *Absence) AbsenceGetAll(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	absences, err := h.absence.FindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	type AbsenceReturn struct {
		ID              uint
		User            model.UserResponse
		AbsenceFrom     time.Time
		AbsenceTill     time.Time
		AbsenceReasonID uint
		NettoDays       int
		CreatedAt       time.Time
	}
	result := []AbsenceReturn{}
	for _, absence := range absences {
		result = append(result, AbsenceReturn{
			ID:              absence.ID,
			User:            user.GetUserResponse(),
			AbsenceFrom:     absence.AbsenceFrom,
			AbsenceTill:     absence.AbsenceTill,
			AbsenceReasonID: *absence.AbsenceReasonID,
			NettoDays:       absence.GetAbsenceWorkDays(),
			CreatedAt:       absence.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Absence) AbsenceReasonsGetAll(c *gin.Context) {
	absenceReasons, err := h.absence.FindAllAbsenceReasons()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(absenceReasons))
}

func (h *Absence) AbsenceCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	var absenceCreateRequest model.AbsenceCreateRequest

	err = c.BindJSON(&absenceCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	absenceReason, err := h.absence.FindAbsenceReasonByID(absenceCreateRequest.AbsenceReasonID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	absence := model.Absence{
		UserID:        &user.ID,
		AbsenceFrom:   absenceCreateRequest.AbsenceFrom,
		AbsenceTill:   absenceCreateRequest.AbsenceTill,
		AbsenceReason: absenceReason,
	}

	err = h.absence.Insert(&absence)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(absence))
}

func (h *Absence) AbsenceQueryUsersSummary(c *gin.Context) {
	absences, err := h.absence.FindByQuery(true, "absence_till >= ?", time.Now().Format("2006-01-02"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	type AbsenceReturn struct {
		ID          uint
		User        model.UserResponse
		AbsenceFrom time.Time
		AbsenceTill time.Time
		Reason      string
		NettoDays   int
	}

	result := []AbsenceReturn{}

	for _, absence := range absences {
		absenceReturn := AbsenceReturn{
			ID:          absence.ID,
			User:        absence.User.GetUserResponse(),
			AbsenceFrom: absence.AbsenceFrom,
			AbsenceTill: absence.AbsenceTill,
			NettoDays:   absence.GetAbsenceWorkDays(),
		}

		if auth.IsAdministrator(c) {
			absenceReturn.Reason = absence.AbsenceReason.Description
		} else {
			absenceReturn.Reason = "Abwesend"
		}

		result = append(result, absenceReturn)
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Absence) AbsenceQueryCurrentUserSummary(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	absences, err := h.absence.FindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	result := model.AbsenceUserSummary{
		ByYear:             make(map[int]model.AbsenceUserSummaryYear),
		HolidayDaysPerYear: user.HolidayDaysPerYear,
	}

	for _, absence := range absences {
		absenceYear := absence.AbsenceFrom.Year()
		if _, exists := result.ByYear[absenceYear]; !exists {
			result.ByYear[absenceYear] = model.AbsenceUserSummaryYear{
				ByAbsenceReason: make(map[uint]model.AbsenceUserSummaryYearReason),
			}
		}

		if _, exists := result.ByYear[absenceYear].ByAbsenceReason[*absence.AbsenceReasonID]; !exists {
			result.ByYear[absenceYear].ByAbsenceReason[*absence.AbsenceReasonID] = model.AbsenceUserSummaryYearReason{}
		}

		yearReasonSummary := result.ByYear[absenceYear].ByAbsenceReason[*absence.AbsenceReasonID]
		days := absence.GetAbsenceWorkDays()

		if absence.AbsenceFrom.Before(time.Now()) {
			yearReasonSummary.Past += int(days)
		} else {
			yearReasonSummary.Upcoming += int(days)
		}

		result.ByYear[absenceYear].ByAbsenceReason[*absence.AbsenceReasonID] = yearReasonSummary
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}
