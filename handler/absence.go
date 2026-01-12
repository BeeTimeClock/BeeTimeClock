package handler

import (
	"errors"
	"fmt"
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

type Absence struct {
	env     *core.Environment
	user    *repository.User
	team    *repository.Team
	absence *repository.Absence
	holiday *repository.Holiday
}

func NewAbsence(env *core.Environment, user *repository.User, absence *repository.Absence, team *repository.Team, holiday *repository.Holiday) *Absence {
	return &Absence{
		env:     env,
		user:    user,
		absence: absence,
		team:    team,
		holiday: holiday,
	}
}

func (h *Absence) AdministrationAbsenceReasonCreate(c *gin.Context) {
	var absenceReasonCreateRequest model.AbsenceReasonCreateRequest
	err := c.BindJSON(&absenceReasonCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	absenceReason := model.AbsenceReason{
		Description:    absenceReasonCreateRequest.Description,
		OvertimeImpact: absenceReasonCreateRequest.OvertimeImpact,
		ImpactHours:    absenceReasonCreateRequest.ImpactDays,
		ImpactDays:     absenceReasonCreateRequest.ImpactHours,
		NeedsApproval:  absenceReasonCreateRequest.NeedsApproval,
	}

	err = h.absence.InsertAbsenceReason(&absenceReason)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(absenceReason))
}

func (h *Absence) AdministrationAbsenceReasonDelete(c *gin.Context) {
	absenceReasonIdParam := c.Param("absenceReasonID")
	absenceReasonId, err := strconv.Atoi(absenceReasonIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	absenceReason, err := h.absence.FindAbsenceReasonByID(uint(absenceReasonId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = h.absence.DeleteAbsenceReason(&absenceReason)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Absence) AdministrationAbsenceReasonUpdate(c *gin.Context) {
	var absenceReasonUpdateRequest model.AbsenceReasonCreateRequest
	err := c.BindJSON(&absenceReasonUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	absenceReasonIdParam := c.Param("absenceReasonID")
	absenceReasonId, err := strconv.Atoi(absenceReasonIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	absenceReason, err := h.absence.FindAbsenceReasonByID(uint(absenceReasonId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	absenceReason.Description = absenceReasonUpdateRequest.Description
	absenceReason.NeedsApproval = absenceReasonUpdateRequest.NeedsApproval
	absenceReason.ImpactDays = absenceReasonUpdateRequest.ImpactDays
	absenceReason.ImpactHours = absenceReasonUpdateRequest.ImpactHours
	absenceReason.OvertimeImpact = absenceReasonUpdateRequest.OvertimeImpact

	err = h.absence.UpdateAbsenceReason(&absenceReason)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(absenceReason))
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

	result := model.AbsenceReturns(absences, &user, false, false, true)

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

	absence, success := h.absenceCreate(c, &user, &absenceCreateRequest, nil)
	if !success {
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(absence))
}

func (h *Absence) absenceCreate(c *gin.Context, user *model.User, absenceCreateRequest *model.AbsenceCreateRequest, signingUser *model.User) (*model.Absence, bool) {
	absenceReason, err := h.absence.FindAbsenceReasonByID(absenceCreateRequest.AbsenceReasonID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return nil, false
	}

	absenceFrom, err := absenceCreateRequest.AbsenceFromParsed()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return nil, false
	}

	absenceTill, err := absenceCreateRequest.AbsenceTillParsed()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return nil, false
	}

	absence := model.Absence{
		UserID:        &user.ID,
		AbsenceFrom:   absenceFrom,
		AbsenceTill:   absenceTill,
		AbsenceReason: absenceReason,
		Identifier:    uuid.New(),
	}

	if absenceReason.NeedsApproval == nil || *absenceReason.NeedsApproval == false {
		absence.Sign(user, model.SIGNED_STATUS_ACCEPTED, nil)
	} else {
		if signingUser != nil {
			absence.Sign(signingUser, model.SIGNED_STATUS_ACCEPTED, nil)
		}
	}

	holidays, err := h.holiday.HolidayFindByDateRange(absenceFrom, absenceTill)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return nil, false
	}

	absence.CalculateNettoDays(holidays)

	err = h.absence.Insert(&absence)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return nil, false
	}

	if microsoft.IsMicrosoftConnected() {
		eventId, err := microsoft.CreateCalendarEntryFromAbsence(user.Username, &absence)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return nil, false
		}

		absenceExternalEvent := model.AbsenceExternalEvent{
			Absence:               absence,
			ExternalEventID:       eventId,
			ExternalEventProvider: model.EXTERNAL_EVENT_PROVIDER_MICROSOFT,
		}

		err = h.absence.AbsenceExternalEventInsert(&absenceExternalEvent)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return nil, false
		}
	}

	return &absence, true
}

func (h *Absence) AbsenceDelete(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	absenceIdParam := c.Param("id")
	absenceId, err := strconv.Atoi(absenceIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	absence, err := h.absence.FindByID(uint(absenceId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if user.AccessLevel == model.USER_ACCESS_LEVEL_USER {
		if *absence.UserID != user.ID {
			c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(fmt.Errorf("not your own absence, can't delete")))
			return
		}

		if !absence.IsDeletableByUser() {
			c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(fmt.Errorf("can't delete past absence")))
			return
		}
	}

	if len(absence.ExternalEvents) > 0 {
		for _, externalEvent := range absence.ExternalEvents {
			if externalEvent.ExternalEventProvider == model.EXTERNAL_EVENT_PROVIDER_MICROSOFT {
				err = microsoft.DeleteCalendarEntry(absence.User.Username, externalEvent.ExternalEventID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
					return
				}
			}
			fmt.Println("delete")
			err = h.absence.AbsenceExternalEventDelete(&externalEvent)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
				return
			}
		}
	}

	err = h.absence.Delete(&absence)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Absence) AbsenceQueryUserYear(c *gin.Context) {
	user, success := getUserFromParam(c, h.user)
	if !success {
		return
	}

	year, success := getYearFromParam(c)
	if !success {
		return
	}

	absences, err := h.absence.FindByUserIDAndYear(user.ID, year)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(model.AbsenceReturns(absences, &user, true, true, true)))
}

func (h *Absence) AbsenceQueryUserYears(c *gin.Context) {
	user, success := getUserFromParam(c, h.user)
	if !success {
		return
	}

	years, err := h.absence.FindYearsWithAbsencesByUserId(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(years))
}

func (h *Absence) AbsenceQueryUsersSummary(c *gin.Context) {
	absences, err := h.absence.FindByQuery(true, "absence_till >= ?", time.Now().Format("2006-01-02"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	result := model.AbsenceReturns(absences, nil, true, auth.IsAdministrator(c), false)
	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Absence) TeamUserAbsenceCreate(c *gin.Context) {
	var absenceCreateRequest model.AbsenceCreateRequest
	err := c.BindJSON(&absenceCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	userId, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	user, err := h.user.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	team, success := getTeamFromParam(c, h.team)
	if !success {
		return
	}

	executingUser, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	isLead := slices.ContainsFunc(team.Members, func(member model.TeamMember) bool {
		return member.UserID == executingUser.ID && (member.Level == model.TeamLevel_Lead || member.Level == model.TeamLevel_LeadSurrogate)
	})

	if !isLead {
		c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(errors.New("you're not lead of the team")))
		return
	}

	if !slices.ContainsFunc(team.Members, func(member model.TeamMember) bool {
		return member.UserID == user.ID
	}) {
		c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(errors.New("user is not member of team")))
		return
	}

	absence, success := h.absenceCreate(c, &user, &absenceCreateRequest, &executingUser)
	if !success {
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(absence))
}

func (h *Absence) AbsenceQueryTeamUsersSummary(c *gin.Context) {
	team, success := getTeamFromParam(c, h.team)
	if !success {
		return
	}

	teamMemberIds := []uint{}
	for _, member := range team.Members {
		teamMemberIds = append(teamMemberIds, member.UserID)
	}

	absences, err := h.absence.FindByQuery(true, "user_id in (?) and absence_till >= ? is not null", teamMemberIds, time.Now().Format("2006-01-02"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	result := model.AbsenceReturns(absences, nil, true, auth.IsAdministrator(c), false)
	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Absence) AbsenceTeamOpen(c *gin.Context) {
	team, success := getTeamFromParam(c, h.team)
	if !success {
		return
	}

	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	isLead := slices.ContainsFunc(team.Members, func(member model.TeamMember) bool {
		return member.UserID == user.ID && (member.Level == model.TeamLevel_Lead || member.Level == model.TeamLevel_LeadSurrogate)
	})

	if !isLead {
		c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(errors.New("you're not lead of the team")))
		return
	}

	userIds := []uint{}

	for _, member := range team.Members {
		userIds = append(userIds, member.UserID)
	}

	abseneceReasons, err := h.absence.FindAllAbsenceReasons()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	absenceReasonIds := []uint{}
	for _, absenceReason := range abseneceReasons {
		if absenceReason.NeedsApproval != nil && *absenceReason.NeedsApproval {
			absenceReasonIds = append(absenceReasonIds, absenceReason.ID)
		}
	}

	teamAbsences := []model.Absence{}
	if len(absenceReasonIds) > 0 {
		teamAbsences, err = h.absence.FindByQuery(true, "user_id in ? and absence_reason_id in ? and signed_user_id is null", userIds, absenceReasonIds)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}
	}

	groupedByUser := make(map[uint][]model.Absence)

	for _, teamAbsence := range teamAbsences {
		userId := *teamAbsence.UserID
		if _, exists := groupedByUser[userId]; !exists {
			groupedByUser[userId] = []model.Absence{}
		}

		groupedByUser[userId] = append(groupedByUser[userId], teamAbsence)
	}

	result := []model.AbsenceReturn{}

	for _, absences := range groupedByUser {
		result = append(result, model.AbsenceReturns(absences, absences[0].User, true, true, false)...)
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Absence) AbsenceQueryUsersSummaryCurrentYear(c *gin.Context) {
	absences, err := h.absence.FindByQuery(true, "absence_till >= ? and absence_till <= ?",
		fmt.Sprintf("%d-01-01", time.Now().Year()), fmt.Sprintf("%d-12-31", time.Now().Year()))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	result := model.AbsenceReturns(absences, nil, true, auth.IsAdministrator(c), false)
	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Absence) AbsenceQueryUsersSummaryCurrentWeek(c *gin.Context) {
	now := time.Now()
	year, week := now.ISOWeek()

	weekStart := helper.WeekStart(year, week)
	weekEnd := weekStart.AddDate(0, 0, 5)

	absences, err := h.absence.FindByQuery(true, "absence_till >= ? and absence_till <= ?",
		weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	result := model.AbsenceReturns(absences, nil, true, auth.IsAdministrator(c), false)
	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *Absence) AbsenceQueryUserSummaryYear(c *gin.Context) {
	userIdParam := c.Param("userID")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	user, err := h.user.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	absences, err := h.absence.FindByUserIDAndYear(user.ID, year)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse(err))
		return
	}

	summary := h.summaryAbsences(&user, absences)

	c.JSON(http.StatusOK, model.NewSuccessResponse(summary))
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

	summary := h.summaryAbsences(&user, absences)

	c.JSON(http.StatusOK, model.NewSuccessResponse(summary))
}

func (h *Absence) summaryAbsences(user *model.User, absences []model.Absence) model.AbsenceUserSummary {
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
		days := *absence.NettoDays

		if absence.AbsenceFrom.Before(time.Now()) {
			yearReasonSummary.Past += days
		} else {
			yearReasonSummary.Upcoming += days
		}

		result.ByYear[absenceYear].ByAbsenceReason[*absence.AbsenceReasonID] = yearReasonSummary
	}

	return result
}

func (h *Absence) AbsenceRecalculate(c *gin.Context) {
	absences, err := h.absence.FindAll(false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	safedHoliays := make(map[int]model.Holidays)
	getHolidays := func(year int) (model.Holidays, error) {
		if value, exists := safedHoliays[year]; exists {
			return value, nil
		}

		holidays, err := h.holiday.HolidayFindByYear(year)
		return holidays, err
	}

	for _, absence := range absences {
		currentNetto := absence.NettoDays

		holidays, err := getHolidays(absence.AbsenceFrom.Year())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
			return
		}
		absence.CalculateNettoDays(holidays)

		if absence.NettoDays != currentNetto {
			err = h.absence.Update(&absence)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
				return
			}
		}
	}
}

func (h *Absence) AbsenceSign(c *gin.Context) {
	var absenceSignRequest model.AbsenceSignRequest
	err := c.BindJSON(&absenceSignRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	team, success := getTeamFromParam(c, h.team)
	if !success {
		return
	}

	executingUser, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	isLead := slices.ContainsFunc(team.Members, func(member model.TeamMember) bool {
		return member.UserID == executingUser.ID && (member.Level == model.TeamLevel_Lead || member.Level == model.TeamLevel_LeadSurrogate)
	})

	if !isLead {
		c.AbortWithStatusJSON(http.StatusForbidden, model.NewErrorResponse(errors.New("you're not lead of the team")))
		return
	}

	absenceId, err := strconv.Atoi(c.Param("absenceID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	absence, err := h.absence.FindByID(uint(absenceId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	absence.Sign(&executingUser, absenceSignRequest.Status, absenceSignRequest.Message)

	err = h.absence.Update(&absence)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(absence))
}
