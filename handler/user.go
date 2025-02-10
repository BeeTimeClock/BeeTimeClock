package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BeeTimeClock/BeeTimeClock-Server/auth"
	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/BeeTimeClock/BeeTimeClock-Server/helper"
	"github.com/BeeTimeClock/BeeTimeClock-Server/model"
	"github.com/BeeTimeClock/BeeTimeClock-Server/repository"
	"github.com/gin-gonic/gin"
)

type User struct {
	env  *core.Environment
	user *repository.User
	team *repository.Team
}

func NewUser(env *core.Environment, user *repository.User, team *repository.Team) *User {
	return &User{
		env:  env,
		user: user,
		team: team,
	}
}

func (h *User) AdministrationUserGetAll(c *gin.Context) {
	withDataQueryParam := c.Query("with_data")
	withData := withDataQueryParam == "true"

	users, err := h.user.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	fmt.Println(withData)

	var result []model.UserResponse
	for _, user := range users {
		result = append(result, user.GetUserResponse())
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *User) AdministrationUserGetByUserID(c *gin.Context) {
	userIdParam := c.Param("userID")
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing userid")))
		return
	}

	user, err := h.user.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.GetUserResponse()))
}

func (h *User) AdministrationUserUpdate(c *gin.Context) {
	var userUpdateRequest model.UserUpdateRequest
	err := c.BindJSON(&userUpdateRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	userIdParam := c.Param("userID")
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing userid")))
		return
	}

	user, err := h.user.FindByID(uint(userId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	user.FirstName = userUpdateRequest.FirstName
	user.LastName = userUpdateRequest.LastName
	user.AccessLevel = userUpdateRequest.AccessLevel
	user.OvertimeSubtractionAmount = userUpdateRequest.OvertimeSubtractionAmount
	user.OvertimeSubtractionModel = userUpdateRequest.OvertimeSubtractionModel

	err = h.user.Update(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.GetUserResponse()))
}

func (h *User) AdministrationUserCreate(c *gin.Context) {
	var userCreateRequest model.UserCreateRequest
	err := c.BindJSON(&userCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	user := model.NewUser(userCreateRequest.Username)
	user.AccessLevel = userCreateRequest.AccessLevel
	user.FirstName = userCreateRequest.FirstName
	user.LastName = userCreateRequest.LastName

	err = user.SetPassword(userCreateRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = h.user.Insert(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(user.GetUserResponse()))
}

func (h *User) AdministrationUserDelete(c *gin.Context) {
	userIDParam := c.Param("userID")

	if strings.TrimSpace(userIDParam) == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("userID missing")))
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := h.user.FindByID(uint(userID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, err)
		return
	}

	err = h.user.Delete(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *User) CurrentUserGet(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(user.GetUserResponse()))
}

func (h *User) CurrentUserApikeyGet(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	apikeys, err := h.user.UserApikeyFindAllByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	result := []model.UserApikeyResponse{}
	for _, apikey := range apikeys {
		result = append(result, apikey.GetUserApikeyResponse())
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(result))
}

func (h *User) CurrentUserApikeyCreate(c *gin.Context) {
	user, err := auth.GetUserFromSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	var userApikeyCreateRequest model.UserApikeyCreateRequest
	err = c.BindJSON(&userApikeyCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	userApikey := model.UserApikey{
		Description: userApikeyCreateRequest.Description,
		User:        user,
		Apikey:      helper.RandomString(64),
		ValidTill:   userApikeyCreateRequest.ValidTill,
	}

	err = h.user.UserApikeyInsert(&userApikey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(userApikey))
}

func (h *User) AdministrationTeamGetAll(c *gin.Context) {
	withDataQueryParam := c.Query("with_data")
	withData := withDataQueryParam == "true"

	teams, err := h.team.TeamFindAll(withData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(teams))
}

func (h *User) AdministrationTeamCreate(c *gin.Context) {
	var teamCreateRequest model.TeamCreateRequest

	err := c.BindJSON(&teamCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	teamOwner, err := h.user.FindByID(teamCreateRequest.TeamOwnerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	team := model.Team{
		Teamname:  teamCreateRequest.Teamname,
		TeamOwner: teamOwner,
	}

	err = h.team.TeamInsert(&team)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(team))
}

func (h *User) AdministrationTeamUpdate(c *gin.Context) {
	var teamUpdateRequest model.TeamCreateRequest
	err := c.BindJSON(&teamUpdateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	_, err = h.user.FindByID(teamUpdateRequest.TeamOwnerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	teamIdParam := c.Param("teamID")
	teamId, err := strconv.ParseUint(teamIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing teamid")))
		return
	}

	team, err := h.team.TeamFindById(uint(teamId), false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	team.Teamname = teamUpdateRequest.Teamname
	team.TeamOwnerID = teamUpdateRequest.TeamOwnerID

	err = h.team.TeamUpdate(&team)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(team))
}

func (h *User) AdministrationTeamGetByID(c *gin.Context) {
	withDataQueryParam := c.Query("with_data")
	withData := withDataQueryParam == "true"

	teamIdParam := c.Param("teamID")
	teamId, err := strconv.ParseUint(teamIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing teamid")))
		return
	}

	team, err := h.team.TeamFindById(uint(teamId), withData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(team))
}

func (h *User) AdministrationTeamMemberGetByTeamID(c *gin.Context) {
	withDataQueryParam := c.Query("with_data")
	withData := withDataQueryParam == "true"

	teamIdParam := c.Param("teamID")
	teamId, err := strconv.ParseUint(teamIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing teamid")))
		return
	}

	members, err := h.team.TeamMemberFindByTeamId(uint(teamId), withData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(members))
}

func (h *User) AdministrationTeamDelete(c *gin.Context) {
	teamIdParam := c.Param("teamID")
	teamId, err := strconv.ParseUint(teamIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing teamid")))
		return
	}

	team, err := h.team.TeamFindById(uint(teamId), false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	err = h.team.TeamDelete(&team)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *User) AdministrationTeamMemberCreate(c *gin.Context) {
	var teamMemberCreateRequest model.TeamMemberCreateRequest
	err := c.BindJSON(&teamMemberCreateRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(err))
		return
	}

	teamIdParam := c.Param("teamID")
	teamId, err := strconv.ParseUint(teamIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing teamid")))
		return
	}

	team, err := h.team.TeamFindById(uint(teamId), true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	user, err := h.user.FindByID(teamMemberCreateRequest.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	teamMember := model.TeamMember{
		Team: team,
		User: user,
	}

	err = h.team.TeamMemberInsert(&teamMember)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse(teamMember))
}

func (h *User) AdministrationTeamMemberDelete(c *gin.Context) {
	teamIdParam := c.Param("teamID")
	teamId, err := strconv.ParseUint(teamIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing teamid")))
		return
	}

	teamMemberIdParam := c.Param("teamMemberID")
	teamMemberId, err := strconv.ParseUint(teamMemberIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.NewErrorResponse(fmt.Errorf("missing teamMemberid")))
		return
	}

	teamMember, err := h.team.TeamMemberFindById(uint(teamMemberId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	if teamMember.TeamID != uint(teamId) {
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Errorf("member is not part of the team"))
		return
	}

	err = h.team.TeamMemberDelete(&teamMember)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.NewErrorResponse(err))
		return
	}

	c.Status(http.StatusNoContent)
}
