package handler

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

func checkUserIsUserTeamlead(c *gin.Context, team *model.Team, executingUser *model.User, teamMember *model.User) (bool, error) {
	isLead := slices.ContainsFunc(team.Members, func(member model.TeamMember) bool {
		return member.UserID == executingUser.ID && (member.Level == model.TeamLevel_Lead || member.Level == model.TeamLevel_LeadSurrogate)
	})

	if !isLead {
		return false, errors.New("you're not lead of the team")
	}

	if !slices.ContainsFunc(team.Members, func(member model.TeamMember) bool {
		return member.UserID == teamMember.ID
	}) {
		return false, errors.New("user is not member of team")
	}

	return true, nil
}

func getTimestampFromParam(c *gin.Context, timestampRepo *repository.Timestamp, userId *uint) (model.Timestamp, bool) {
	timestampIdParam := c.Param("timestampID")
	timestampId, err := strconv.Atoi(timestampIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return model.Timestamp{}, false
	}

	timestamp, err := timestampRepo.FindByID(uint(timestampId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return model.Timestamp{}, false
	}

	if userId != nil {
		if timestamp.UserID != *userId {
			c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("user (%d) mismatch for timestamp (%d)", userId, timestamp.ID)))
			return model.Timestamp{}, false
		}
	}

	return timestamp, true
}

func getUserFromParam(c *gin.Context, userRepo *repository.User) (model.User, bool) {
	userIdParam := c.Param("userID")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return model.User{}, false
	}

	user, err := userRepo.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return model.User{}, false
	}

	return user, true
}

func getTeamFromParam(c *gin.Context, teamRepo *repository.Team) (model.Team, bool) {
	teamIdParam := c.Param("teamID")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return model.Team{}, false
	}

	team, err := teamRepo.TeamFindById(uint(teamId), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return model.Team{}, false
	}

	return team, true
}

func getYearFromParam(c *gin.Context) (int, bool) {
	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return 0, false
	}

	return year, true
}

func getYearMonthFromParam(c *gin.Context) (int, int, bool) {
	year, success := getYearFromParam(c)
	if !success {
		return year, 0, false
	}

	monthParam := c.Param("month")
	month, err := strconv.Atoi(monthParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return year, 0, false
	}

	return year, month, true
}

func getClientIPByHeaders(c *gin.Context) (ip string, err error) {
	headers := []string{
		"X-Forwarded-For",
		"x-forwarded-for",
		"X-FORWARDED-FOR",
		"X-Real-Ip",
	}

	for _, header := range headers {
		value := c.Request.Header.Get(header)
		if value != "" {
			return value, nil
		}
	}

	return "", errors.New("cant detect client ip")
}
